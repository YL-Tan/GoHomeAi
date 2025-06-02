# import httpx
# import pytest

# @pytest.mark.integration
# def test_healthcheck():
#     """
#     Requires ML API container running locally on port 9999.
#     The GitHub workflow will spin it up via docker-compose later.
#     """
#     resp = httpx.get("http://localhost:9999/health")
#     assert resp.status_code == 200
#     assert resp.json() == {"status": "ok"}
