from typing import Optional, Tuple, List
from vyper import v
from dataclasses import dataclass
from nats.js.kv import KeyValue
from nats.js.client import JetStreamContext
from loguru import logger
from loguru._logger import Logger
from nats.js.errors import KeyNotFoundError

from exceptions import (
    FailedGettingConfigGivenKey,
    FailedInitializingConfig,
    FailedGettingConfig,
    FailedSettingConfig,
    FailedDeletingConfig,
)

from constants import Scope

@dataclass
class CentralizedConfig:
    product_kv: KeyValue
    workflow_kv: KeyValue
    process_kv: KeyValue
    js: JetStreamContext
    logger: Logger = logger.bind(component="[CENTRALIZED CONFIGURATION]")

    async def initialize(self):
        self.product_kv, self.workflow_kv, self.process_kv = await self.init_kv_stores()
    
    async def init_kv_stores(self) -> Tuple[KeyValue, KeyValue, KeyValue]:
        try:
            name = v.get("centralized_configuration.product.bucket")
            logger.info(f"initializing product key-value store {name}...")
            product_kv = await self.js.key_value(bucket=name)
            logger.info("product key-value store initialized!")

            name = v.get("centralized_configuration.workflow.bucket")
            logger.info(f"initializing workflow key-value store {name}...")
            workflow_kv = await self.js.key_value(bucket=name)
            logger.info("workflow key-value store initialized!")

            name = v.get("centralized_configuration.process.bucket")
            logger.info(f"initializing process key-value store {name}...")
            process_kv = await self.js.key_value(bucket=name)
            logger.info("process key-value store initialized!")

            return product_kv, workflow_kv, process_kv
        except Exception as e:
            raise FailedInitializingConfig(error=e)

    def get_config(self, key: str, scope: str = "") -> str:
            if scope:
                try:
                    config = self.get_config_from_scope(key, scope)
                except KeyNotFoundError:
                    raise FailedGettingConfig(key=key, scope=scope, error=KeyNotFoundError)
                return config

            all_scopes_in_order = [Scope.ProcessScope, Scope.WorkflowScope, Scope.ProductScope]
            for _scope in all_scopes_in_order:
                try:
                    config = self.get_config_from_scope(key, _scope)
                except KeyNotFoundError:
                    continue    
                except Exception as e:
                    raise FailedGettingConfig(key=key, scope=_scope, error=e)

                return config

            raise FailedGettingConfigGivenKey(key=key)

    def set_config(self, key: str, value: str, scope: str = ""):
        kv_store = self.get_scoped_config(scope)

        try:
            kv_store.PutString(key, value)
        except Exception as e:
            raise FailedSettingConfig(key=key, scope=scope, error=e)

    def delete_config(self, key: str, scope: str=""):
        try:
            self.get_scoped_config(scope).delete(key)
        except Exception as e:
            raise FailedDeletingConfig(key=key, scope=scope, error=e)

    def get_config_from_scope(self, key: str, scope: str="") -> Tuple[str, Optional[Exception]]:
        return self.get_scoped_config(scope).get(key).value().decode('utf-8')

    def get_scoped_config(self, scope: str="") -> KeyValue:
        if not scope or Scope(scope) == Scope.ProcessScope:
            return self.process_kv
        elif Scope(scope) == Scope.ProductScope:
            return self.product_kv
        elif Scope(scope) == Scope.WorkflowScope:
            return self.workflow_kv
