import glob
import os
import sys
from pdfrw import PdfReader, PdfWriter


def error(*args, **kwargs):
    kwargs["file"] = sys.stderr
    print(*args, **kwargs)


def fatal(*args, **kwargs):
    error(*args, **kwargs)
    exit(1)


if __name__ == "__main__":
    if len(sys.argv) < 3:
        fatal("Usage:", sys.argv[0], "FILE [FILE ...]")

    cfg = {"out": sys.argv[1], "files": []}

    for pattern in sys.argv[2:]:
        for path in glob.iglob(pattern):
            cfg["files"].append(path)

    wr = PdfWriter()
    for path in cfg["files"]:
        wr.addpages(PdfReader(path).pages)

    wr.write(cfg["out"])

