from dataclasses import dataclass
from typing import Optional

from centralized_config.constants import Scope
from centralized_config.exceptions import FailedDeletingConfig, FailedGettingConfig, FailedInitializingConfig, FailedSettingConfig
from loguru import logger
from loguru._logger import Logger
from nats.js.client import JetStreamContext
from nats.js.errors import KeyNotFoundError
from nats.js.kv import KeyValue
from vyper import v


@dataclass
class CentralizedConfig:
    product_kv: KeyValue
    workflow_kv: KeyValue
    process_kv: KeyValue
    js: JetStreamContext
    logger: Logger = logger.bind(component="[CENTRALIZED CONFIGURATION]")

    async def initialize(self) -> Optional[Exception]:
        self.product_kv, self.workflow_kv, self.process_kv = await self._init_kv_stores()

    async def _init_kv_stores(self) -> tuple[KeyValue, KeyValue, KeyValue] | Exception:
        try:
            name = v.get("centralized_configuration.product.bucket")
            logger.info(f"initializing product key-value store {name}...")
            product_kv = await self.js.key_value(bucket=name)
            logger.info("product key-value store initialized")

            name = v.get("centralized_configuration.workflow.bucket")
            logger.info(f"initializing workflow key-value store {name}...")
            workflow_kv = await self.js.key_value(bucket=name)
            logger.info("workflow key-value store initialized")

            name = v.get("centralized_configuration.process.bucket")
            logger.info(f"initializing process key-value store {name}...")
            process_kv = await self.js.key_value(bucket=name)
            logger.info("process key-value store initialized")

            return product_kv, workflow_kv, process_kv
        except Exception as e:
            raise FailedInitializingConfig(error=e)

    def get_config(self, key: str, scope: Optional[Scope] = None) -> tuple[str, bool] | Exception:
        if scope:
            try:
                config = self._get_config_from_scope(key, scope)
            except KeyNotFoundError:
                logger.warning(f"key {key} not found in scope {scope}")
                return None, False
            except Exception as e:
                raise FailedGettingConfig(key=key, scope=scope, error=e)

            return config, True

        for _scope in Scope:
            try:
                config = self._get_config_from_scope(key, _scope)
            except KeyNotFoundError:
                logger.debug(f"key {key} not found in scope {_scope}")
                continue
            except Exception as e:
                raise FailedGettingConfig(key=key, scope=_scope, error=e)

            return config, True

        logger.warning(f"key {key} not found in any scope")
        return None, False

    def set_config(self, key: str, value: str, scope: Optional[Scope] = None) -> Optional[Exception]:
        kv_store = self._get_scoped_config(scope)

        try:
            kv_store.Put(key, value)
        except Exception as e:
            raise FailedSettingConfig(key=key, scope=scope, error=e)

    def delete_config(self, key: str, scope: Optional[Scope] = None) -> bool | Exception:
        try:
            return self._get_scoped_config(scope).delete(key)
        except Exception as e:
            raise FailedDeletingConfig(key=key, scope=scope, error=e)

    def _get_config_from_scope(self, key: str, scope: str) -> str:
        return self._get_scoped_config(scope).get(key).value().decode("utf-8")

    def _get_scoped_config(self, scope: Optional[Scope] = Scope.ProcessScope) -> KeyValue:
        if scope == Scope.ProductScope:
            return self.product_kv
        elif scope == Scope.WorkflowScope:
            return self.workflow_kv
        elif scope == Scope.ProcessScope:
            return self.process_kv
