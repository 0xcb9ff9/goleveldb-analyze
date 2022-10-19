#!/usr/bin/env python

import argparse

parser = argparse.ArgumentParser(description='show key/value size distribution')

parser.add_argument('db', type=str, help='sqlite database path')

import sqlite3
import pandas as pd

def main():
    args = parser.parse_args()

    with sqlite3.connect(args.db) as cnx:
        df = pd.read_sql_query("SELECT key_size,value_size,sum_size FROM kvsize", cnx)
        pd.set_option('float_format', '{:f} byte'.format)

        print("key size distribution")
        print(df["key_size"].describe())
        print("")

        print("value size distribution")
        print(df["value_size"].describe())
        print("")

        print("key+value size distribution")
        print(df["sum_size"].describe())

if __name__ == "__main__":
    main()

