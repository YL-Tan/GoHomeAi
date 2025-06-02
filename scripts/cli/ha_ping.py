#!/usr/bin/env python
import click
import httpx


@click.command()
@click.option(
    "--url",
    default="http://localhost:8000/ha/status",
    help="The Home-Assistant mock URL to ping.",
)
@click.option("--timeout", default=5.0, help="Request timeout in seconds.")
def ha_ping(url: str, timeout: float):
    """Simple CLI that pings a Home-Assistant API (mock) endpoint and prints the JSON response.

    Args:
        url (str): The endpoint to ping. Defaults to http://localhost:8000/ha/status.
        timeout (float): The maximum time (in seconds) to wait for a response before timing out.

    Returns:
        SystemExit: Exits with code 1 in case of connection or HTTP error.
    """
    try:
        response = httpx.get(url, timeout=timeout)
        response.raise_for_status()
    except httpx.RequestError as exc:
        click.echo(f"Error: Could not reach {url}, Details: {exc}", err=True)
        return SystemExit(1)
    except httpx.HTTPStatusError as exc:
        click.echo(
            f"Error: Received bad status code {exc.response.status_code}", err=True
        )
        return SystemExit(1)

    click.echo("Successful ping! Response JSON")
    click.echo(response.text)


if __name__ == "__main__":
    ha_ping()
