#!/usr/bin/env python3
"""
Script rendering Go benchmark results into plots
"""

import datetime

import pandas as pd

pd.set_option('display.max_rows', 500)
pd.set_option('display.max_columns', 500)
pd.set_option('display.width', 140)


class BenchmarkLine(object):
    description: str
    run_number: int
    runtime: float
    units: str

    def __init__(self, description: str, run_number: int, runtime: float, units: str):
        self.description = description
        self.run_number = run_number
        self.runtime = runtime

        if units != "ns/op":
            raise ValueError("FIXME: new unit")

        self.units = units

    @classmethod
    def from_str(cls, src: str):
        split = src.split()
        if len(split) != 4:
            raise ValueError(f"invalid line format: {src}")

        return cls(split[0], int(split[1]), float(split[2]), split[3])


class BenchmarkCase(object):
    threads: int
    data_source: str
    set_kind: str
    scenario: str
    runtime: datetime.timedelta

    def __init__(self, threads: int, data_source: str, set_kind: str, scenario: str, runtime: datetime.timedelta):
        self.threads = threads
        self.data_source = data_source
        self.set_kind = set_kind
        self.scenario = scenario
        self.runtime = runtime

    @classmethod
    def from_benchmark_line(cls, line: BenchmarkLine):
        split = line.description.split("/")

        if len(split) != 5:
            raise ValueError(f"invalid description format: {line.description}")

        threads = int(split[1].split("_")[0])
        data_source = split[2]
        set_kind = split[3]
        scenario = split[4]

        if line.units == "ns/op":
            runtime = datetime.timedelta(microseconds=line.runtime / 1000)
        else:
            raise ValueError("FIXME: new unit")

        return BenchmarkCase(threads, data_source, set_kind, scenario, runtime)


def parse_report(src_file: str) -> pd.DataFrame:
    with open(src_file) as f:
        src = f.read()

    raw_lines = src.splitlines()
    raw_lines = raw_lines[4:-2]  # drop extra information
    parsed_lines = [BenchmarkLine.from_str(line) for line in raw_lines]
    bench_cases = [BenchmarkCase.from_benchmark_line(line) for line in parsed_lines]
    df = pd.DataFrame([bc.__dict__ for bc in bench_cases])
    print(df)
    return df


def render_plot(df: pd.DataFrame, scenario: str, data_source: str):
    x = df[df["data_source"] == data_source][["threads", "set_kind", "runtime"]]
    pivot = x.pivot(columns=['set_kind'], values=['runtime'], index='threads')

    # axes = pivot.plot(kind="line", title="Sort.{}".format(method_name), logx=True, logy=True)
    axes = pivot.plot(kind="line", title="Scenario: {}".format(scenario))
    axes.set_ylabel("nanoseconds")
    lgd = axes.legend(loc='center right', bbox_to_anchor=(1.5, 0.5))
    axes.figure.savefig("/tmp/" + scenario + ".svg",
                        bbox_extra_artists=(lgd,), bbox_inches='tight')


def main():
    src_file = "report.txt"
    df = parse_report(src_file)
    render_plot(df, "ololo", "ascending_array_input")


if __name__ == "__main__":
    main()
