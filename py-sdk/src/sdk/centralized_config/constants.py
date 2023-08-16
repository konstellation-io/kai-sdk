from enum import Enum


class Scope(Enum):
    ProcessScope = "process"
    WorkflowScope = "workflow"
    ProductScope = "product"
