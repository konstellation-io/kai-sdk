from vyper import v

from sdk.metadata.metadata import Metadata


def test_ok():
    v.set("metadata.product_id", "test_product_id")
    v.set("metadata.workflow_id", "test_workflow_id")
    v.set("metadata.process_id", "test_process_id")
    v.set("metadata.version_id", "test_version_id")
    v.set("nats.object_store", "test_object_store")
    v.set("centralized_configuration.product.bucket", "test_product")
    v.set("centralized_configuration.workflow.bucket", "test_workflow")
    v.set("centralized_configuration.process.bucket", "test_process")

    metadata = Metadata()

    assert metadata is not None
    assert metadata.logger is not None
    assert metadata.get_product() == "test_product_id"
    assert metadata.get_workflow() == "test_workflow_id"
    assert metadata.get_process() == "test_process_id"
    assert metadata.get_version() == "test_version_id"
    assert metadata.get_object_store_name() == "test_object_store"
    assert metadata.get_key_value_store_product_name() == "test_product"
    assert metadata.get_key_value_store_workflow_name() == "test_workflow"
    assert metadata.get_key_value_store_process_name() == "test_process"