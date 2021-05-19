#!/usr/bin/env python3
"""
Script to render Go benchmark results into plots
"""

import datetime

import pandas as pd

pd.set_option('display.max_rows', 500)
pd.set_option('display.max_columns', 500)
pd.set_option('display.width', 140)


class BenchmarkLine(object):
    description: str
    run_number: int
    benchmark_time: float
    units: str

    def __init__(self, description: str, run_number: int, runtime: float, units: str):
        self.description = description
        self.run_number = run_number
        self.benchmark_time = runtime

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
    benchmark_time: pd.Timedelta

    def __init__(self, threads: int, data_source: str, set_kind: str, scenario: str, benchmark_time: datetime.timedelta):
        self.threads = threads
        self.data_source = data_source
        self.set_kind = set_kind
        self.scenario = scenario
        self.benchmark_time = benchmark_time

    @classmethod
    def from_benchmark_line(cls, line: BenchmarkLine):
        split = line.description.split("/")

        if len(split) != 5:
            raise ValueError(f"invalid description format: {line.description}")

        threads = int(split[1].split("_")[0])
        data_source = split[2]
        set_kind = split[3]
        scenario = split[4].split("-")[0]  # crop tail with a number of CPU cores

        if line.units == "ns/op":
            runtime = pd.Timedelta(line.benchmark_time, unit='ns')
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


def render_plots(df: pd.DataFrame):
    scenarios = df['scenario'].unique()
    data_sources = df['data_source'].unique()

    for scenario in scenarios:
        for data_source in data_sources:
            render_plot(df, scenario, data_source)


def render_plot(df: pd.DataFrame, scenario: str, data_source: str):
    crop = df[df["data_source"] == data_source][df["scenario"] == scenario]
    crop = crop[["threads", "set_kind", "benchmark_time"]] # drop extra columns
    print(crop)
    pivot = crop.pivot(columns='set_kind', values='benchmark_time', index='threads')
    print(pivot)
    # pivot = pivot.droplevel(0, axis=1)
    # axes = pivot.plot(kind="line", title="Sort.{}".format(method_name), logx=True, logy=True)
    axes = pivot.plot(
        kind="line",
        title=f"Scenario: {scenario}, data_source: {data_source}",
        logy=True,
        legend="True",
        colormap="Set1",
    )
    axes.set_ylabel("nanoseconds")
    lgd = axes.legend(loc='best')
    filename = f'./report/{scenario}_{data_source}.svg'
    axes.figure.savefig(filename, bbox_extra_artists=(lgd,), bbox_inches='tight')


def main():
    src_file = "report/report.txt"
    df = parse_report(src_file)
    render_plots(df)


if __name__ == "__main__":
    main()
