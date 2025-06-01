import pathlib
import pickle

import pytest
import torch

MODEL_PATH = pathlib.Path(__file__).resolve().parents[1] / "model.pt"


@pytest.mark.skipif(not MODEL_PATH.exists(), reason="model.pt not present yet")
def test_model_file_loads():
    """
    Even a dummy file should be loadable with torch.load.
    Replace with real deserialisation later.
    """
    try:
        _ = torch.load(MODEL_PATH, map_location="cpu")
    except (RuntimeError, pickle.UnpicklingError):
        pytest.fail("model.pt exists but cannot be loaded by torch.load()")
