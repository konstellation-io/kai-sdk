from unittest.mock import patch

import pytest
from vyper import v

from sdk.measurements.exceptions import FailedToInitializeMeasurementsError
from sdk.measurements.measurements import Measurements


@patch.object(Measurements, "_setup_metrics")
def test_ok(_):
    v.set("opentelemetry.endpoint", "localhost:4317")
    v.set("opentelemetry.insecure", True)
    v.set("opentelemetry.timeout", 5)
    v.set("opentelemetry.metrics_interval", 1)
    m_measurements = Measurements()

    assert m_measurements


@patch.object(Measurements, "_setup_metrics", side_effect=Exception("error"))
def test_ko(_):
    v.set("opentelemetry.endpoint", "localhost:4317")
    v.set("opentelemetry.insecure", True)
    v.set("opentelemetry.timeout", 5)
    v.set("opentelemetry.metrics_interval", 1)
    with pytest.raises(FailedToInitializeMeasurementsError):
        measurements = Measurements()

        assert measurements.meter is None
