from click.testing import CliRunner
from ha_ping import ha_ping


def test_ha_ping_success(monkeypatch):
    class DummyResponse:
        status_code = 200
        text = '{"status":"online"}'

        def raise_for_status(self):
            pass

    def dummy_get(url, timeout):
        return DummyResponse()

    # Monkeypatch httpx.get
    import httpx

    monkeypatch.setattr(httpx, "get", dummy_get)

    runner = CliRunner()
    result = runner.invoke(ha_ping, ["--url", "http://example.com"])
    assert result.exit_code == 0
    assert "Successful ping!" in result.output
