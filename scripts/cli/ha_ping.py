#!/usr/bin/env python
import os
import time

import click
import httpx
import mlflow

MLFLOW_TRACKING_URIS = {
    "dev": "mlruns/",
    "staging": "",
    "prod": "",
}


@click.command()
@click.option(
    "--env",
    type=click.Choice(["dev", "staging", "prod"], case_sensitive=False),
    default="dev",
    help="Which environment to run in: dev, staging, or prod.",
)
@click.option(
    "--url",
    default="http://localhost:9999/ha/status",
    help="The Home-Assistant mock URL to ping.",
)
@click.option("--timeout", default=5.0, help="Request timeout in seconds.")
@click.option(
    "--mlflow-experiment-name",
    default="ha_ping_cli",
    help="Name of the MLflow experiment under which to log runs.",
)
def ha_ping(env: str, url: str, timeout: float, mlflow_experiment_name: str):
    """CLI that pings a Home-Assistant API (mock) endpoint, logs via MLflow, and prints the JSON response.

    Args:
        url (str): The endpoint to ping. Defaults to http://localhost:9999/ha/status.
        timeout (float): The maximum time (in seconds) to wait for a response before timing out.
        mlflow_experiment_name (str): The name of the MLflow experiment under which to log runs.

    Returns:
        SystemExit: Exits with code 1 in case of connection or HTTP error.
    """
    # 3. Configure MLflow tracking URI for this environment
    tracking_uri = MLFLOW_TRACKING_URIS[env]
    mlflow.set_tracking_uri(tracking_uri)
    mlflow.set_experiment(mlflow_experiment_name)

    # Make sure tmp/ exists for the artifact
    if env == "dev":
        os.makedirs("tmp", exist_ok=True)

    # Start _one_ run that wraps everything
    with mlflow.start_run(run_name=f"ping_at_{env}_" + time.strftime("%Y%m%d_%H%M%S")):
        mlflow.set_tag("environment", env)
        mlflow.log_param("url", url)
        mlflow.log_param("timeout", timeout)

        start_time = time.time()
        try:
            # Perform the HTTP request inside the same run
            response = httpx.get(url, timeout=timeout)
            response.raise_for_status()
        except httpx.RequestError as exc:
            latency = time.time() - start_time
            mlflow.log_metric("latency_seconds", latency)
            mlflow.log_metric("status_code", 0)
            mlflow.set_tag("success", "false")

            click.echo(f"Error: Could not reach {url}, Details: {exc}", err=True)
            # Exit with code 1 (keeps this run “failed”)
            raise SystemExit(1)

        except httpx.HTTPStatusError as exc:
            latency = time.time() - start_time
            mlflow.log_metric("latency_seconds", latency)
            mlflow.log_metric("status_code", exc.response.status_code)
            mlflow.set_tag("success", "false")

            click.echo(
                f"Error: Received bad status code {exc.response.status_code}", err=True
            )
            raise SystemExit(1)

        # If we reach here, request succeeded
        latency = time.time() - start_time
        status_code = response.status_code

        mlflow.log_metric("latency_seconds", latency)
        mlflow.log_metric("status_code", status_code)
        mlflow.set_tag("success", "true")

        # Save JSON response as an artifact
        if env == "dev":
            run_id = mlflow.active_run().info.run_id
            artifact_path = f"tmp/response_{run_id}.json"
            with open(artifact_path, "w") as f:
                f.write(response.text)
            mlflow.log_artifact(artifact_path)

        click.echo("Successful ping! Response JSON:")
        click.echo(response.text)


if __name__ == "__main__":
    ha_ping()
