def test_core_imports():
    """
    Sanity-check that key libraries import without side-effects.
    """
    import numpy as np
    import torch

    assert np.__version__.startswith(("1.", "1.2")), "Unexpected NumPy major version"
    assert torch.__version__.startswith("2."), "Unexpected PyTorch major version"
