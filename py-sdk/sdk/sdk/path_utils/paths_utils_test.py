from vyper import v

from sdk.path_utils.path_utils import PathUtils


def test_ok():
    PathUtils()


def test_get_basepath_ok():
    v.set("metadata.base_path", "base_path")
    PathUtils()

    base_path = PathUtils.get_base_path()

    assert base_path == "base_path"


def test_compose_path_ok():
    v.set("metadata.base_path", "path")
    PathUtils()

    base_path = PathUtils.compose_path("other_path", "another_path")

    assert base_path == "path/other_path/another_path"
