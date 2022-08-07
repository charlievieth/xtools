#!/usr/bin/env python3

import json
import subprocess
import sys
from typing import Dict


def parse_modfile(name: str) -> Dict[str, str]:
    m = json.loads(subprocess.check_output(
        ["go", "mod", "edit", "-json", name],
        encoding="utf-8",
    ))
    return {v["Path"]: v["Version"] for v in m["Require"]}


def generate_modfile(require: Dict[str, str]) -> str:
    reqs = "\n".join(
        [f"\t{path} {version}" for path, version in require.items()],
    )
    return f"""
module github.com/charlievieth/xtools

go 1.18

require (
    {reqs}
)
    """


if __name__ == "__main__":
    require = {}
    for req in [parse_modfile(name) for name in sys.argv[1:]]:
        require.update(req)
    print(generate_modfile(require))
