from dataclasses import dataclass
from runner.common.common import Handler

Preprocessor = Postprocessor = Handler


@dataclass
class ExitRunner:
    pass
