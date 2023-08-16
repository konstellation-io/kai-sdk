import pytest
from messaging.messaging_utils import (
    FailedGettingMaxMessageSizeError,
    MessagingUtils,
    compress,
    decompress,
    is_compressed,
    size_in_kb,
    size_in_mb,
)
from mock import patch
from nats.aio.client import Client as NatsClient
from nats.js.api import StreamConfig, StreamInfo
from nats.js.client import JetStreamContext
from vyper import v


async def test_ok():
    nc = NatsClient()
    js = nc.jetstream()

    utils = MessagingUtils(js=js, nc=nc)

    assert utils is not None
    assert utils.js is not None
    assert utils.nc is not None


@pytest.mark.parametrize(
    "nats_max_payload, jetstream_max_msg_size, max_message_size",
    [
        (2048, 1024, 1024),
        (1024, 2048, 1024),
        (1024, -1, 1024),
    ],
)
@patch.object(JetStreamContext, "stream_info")
@patch.object(NatsClient, "max_payload")
async def test_get_max_message_size_ok(
    nats_max_payload_mock, jetstream_context_mock, nats_max_payload, jetstream_max_msg_size, max_message_size
):
    v.set("nats.stream", "test_stream")
    nc = NatsClient()
    js = nc.jetstream()
    nats_max_payload_mock.return_value = nats_max_payload
    jetstream_context_mock.return_value = StreamInfo(
        config=StreamConfig(max_msg_size=jetstream_max_msg_size), state=None
    )

    utils = MessagingUtils(js=js, nc=nc)
    max_size = await utils.get_max_message_size()

    assert max_size == max_message_size


@patch.object(JetStreamContext, "stream_info", side_effect=Exception)
@patch.object(NatsClient, "max_payload", return_value=1024)
async def test_get_max_message_size_ko(nats_max_payload_mock, jetstream_context_mock):
    v.set("nats.stream", "test_stream")
    nc = NatsClient()
    js = nc.jetstream()

    with pytest.raises(FailedGettingMaxMessageSizeError):
        utils = MessagingUtils(js=js, nc=nc)
        await utils.get_max_message_size()


def test_size_in_mb_ok():
    assert size_in_mb(1024 * 1024) == "1.00 MB"


def test_size_in_kb_ok():
    assert size_in_kb(1024) == "1.00 KB"


def test_compress_decompress_ok():
    compress_ = compress(b"test")
    decompress_ = decompress(compress_)
    assert decompress_ == b"test"
    assert is_compressed(compress_)
