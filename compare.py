#!/usr/bin/env python3

#
# Reads and compare contents of 2 json files
#
# Usage:
#   compare.py file1.json file2.json
#

import json
import sys

from deepdiff import DeepDiff


def read_dict_from_json_file(filename):
    with open(filename, "r") as fd:
        return json.load(fd)


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage:\n\t{} <file1> <file2>".format(__file__))
    dict0 = read_dict_from_json_file(sys.argv[1])
    dict1 = read_dict_from_json_file(sys.argv[2])

    diff = DeepDiff(dict0, dict1)

    print("These are the keys:", diff.keys())

    print("Asserting there are no removed items")
    assert "dict_items_removed" not in diff.keys()

    print("Asserting there are no type changes")
    assert "type_changes" not in diff.keys()

    print("Asserting only modified key is 'Updated'")
    assert "root['updated']" in diff["values_changed"]

    if "iterable_item_added" in diff:
        print("Items added")
        for k, v in diff["iterable_item_added"].items():
            print("{}: {}".format(k, v))

    if "iterable_item_removed" in diff:
        print("Items removed")
        for k, v in diff["iterable_item_removed"].items():
            print("{}: {}".format(k, v))
