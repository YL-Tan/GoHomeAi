import importlib.metadata as im


def test_numpy_is_v1():
    """
    Guardrail until PyTorch releases wheels linked against NumPy 2.x.
    """
    version = im.version("numpy")
    major = int(version.split(".")[0])
    assert major == 1, f"NumPy {version} is incompatible with current PyTorch"
