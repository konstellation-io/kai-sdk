from dataclasses import dataclass
from typing import Optional

from centralized_config.constants import Scope
from centralized_config.exceptions import (
    FailedDeletingConfigError,
    FailedGettingConfigError,
    FailedInitializingConfigError,
    FailedSettingConfigError,
)
from loguru import logger
from loguru._logger import Logger
from nats.js.client import JetStreamContext
from nats.js.errors import KeyNotFoundError
from nats.js.kv import KeyValue
from vyper import v


@dataclass
class CentralizedConfig:
    js: JetStreamContext
    product_kv: KeyValue = None
    workflow_kv: KeyValue = None
    process_kv: KeyValue = None
    logger: Logger = logger.bind(component="[CENTRALIZED CONFIGURATION]")

    async def initialize(self) -> Optional[Exception]:
        self.product_kv, self.workflow_kv, self.process_kv = await self._init_kv_stores()

    async def _init_kv_stores(self) -> tuple[KeyValue, KeyValue, KeyValue] | Exception:
        try:
            name = v.get("centralized_configuration.product.bucket")
            self.logger.info(f"initializing product key-value store {name}...")
            product_kv = await self.js.key_value(bucket=name)
            self.logger.info("product key-value store initialized")

            name = v.get("centralized_configuration.workflow.bucket")
            self.logger.info(f"initializing workflow key-value store {name}...")
            workflow_kv = await self.js.key_value(bucket=name)
            self.logger.info("workflow key-value store initialized")

            name = v.get("centralized_configuration.process.bucket")
            self.logger.info(f"initializing process key-value store {name}...")
            process_kv = await self.js.key_value(bucket=name)
            self.logger.info("process key-value store initialized")

            return product_kv, workflow_kv, process_kv
        except Exception as e:
            self.logger.warning(f"failed initializing configuration: {e}")
            raise FailedInitializingConfigError(error=e)

    async def get_config(self, key: str, scope: Optional[Scope] = None) -> tuple[str, bool] | Exception:
        if scope:
            try:
                config = await self._get_config_from_scope(key, scope)
            except KeyNotFoundError as e:
                self.logger.debug(f"key {key} not found in scope {scope}: {e}")
                return None, False
            except Exception as e:
                self.logger.warning(f"failed getting config: {e}")
                raise FailedGettingConfigError(key=key, scope=scope, error=e)

            return config, True

        for _scope in Scope:
            try:
                config = await self._get_config_from_scope(key, _scope)
            except KeyNotFoundError as e:
                self.logger.debug(f"key {key} not found in scope {_scope}: {e}")
                continue
            except Exception as e:
                self.logger.warning(f"failed getting config: {e}")
                raise FailedGettingConfigError(key=key, scope=_scope, error=e)

            return config, True

        self.logger.warning(f"key {key} not found in any scope")
        return None, False

    async def set_config(self, key: str, value: bytes, scope: Optional[Scope] = None) -> Optional[Exception]:
        kv_store = self._get_scoped_config(scope)

        try:
            await kv_store.put(key, value)
        except Exception as e:
            self.logger.warning(f"failed setting config: {e}")
            raise FailedSettingConfigError(key=key, scope=scope, error=e)

    async def delete_config(self, key: str, scope: Optional[Scope] = None) -> bool | Exception:
        kv_store = self._get_scoped_config(scope)
        try:
            return await kv_store.delete(key)
        except Exception as e:
            self.logger.warning(f"failed deleting config: {e}")
            raise FailedDeletingConfigError(key=key, scope=scope, error=e)

    async def _get_config_from_scope(self, key: str, scope: Scope) -> str:
        scoped_config = self._get_scoped_config(scope)
        entry = await scoped_config.get(key)
        return entry.value.decode("utf-8")

    def _get_scoped_config(self, scope: Optional[Scope] = Scope.ProcessScope) -> KeyValue:
        if scope == Scope.ProductScope:
            return self.product_kv
        elif scope == Scope.WorkflowScope:
            return self.workflow_kv
        elif scope == Scope.ProcessScope:
            return self.process_kv
