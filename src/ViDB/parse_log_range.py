#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import re
import os

import numpy
import argparse
import pandas as pd
import matplotlib.pyplot as plt

def actual_size(size):
    if "K" in size:
        size_f = float(size.replace("K", "")) * 1000
    elif "M" in size:
        size_f = float(size.replace("M", "")) * 1000 * 1000
    elif "G" in size:
        size_f = float(size.replace("G", "")) * 1000 * 1000 * 1000
    else:
        size_f = float(size)
    return size_f
        
def parse_log(log_file):
    meta_info = log_file.split('/')[-1].split('.')[0].split('_')
    type = meta_info[0].strip()
    entry_count = int(actual_size(meta_info[1].strip()))
    range_size = int(meta_info[2].strip())
    timestamp = meta_info[-2].strip() + meta_info[-1].strip()
    if type == "range-get" or type == "range-query":
        lat_ns, lat_s, rs, tps = parse_range_get_latencies(log_file)
        return type, entry_count, rs, lat_s, tps, timestamp
    else:
        return "None", 0, 0, 0, 0, timestamp
    
def parse_range_get_latencies(log_file):
    latencies_ns = []
    latencies_s = []
    range_sizes = []
    throughputs = []
    with open(log_file, 'r') as f:
        for line in f:
            pattern = r"Latency:\s*(\d+)\s*ns\s*\(([\d.]+)\s*s\),\s*range size:\s*(\d+),\s*TPS:\s*([\d.]+)"
            match = re.search(pattern, line)

            if match:
                latencies_ns.append(int(match.group(1)))
                latencies_s.append(float(match.group(2)))
                range_sizes.append(int(match.group(3)))
                throughputs.append(float(match.group(4)))
            else:
                pass
    return latencies_ns, latencies_s, range_sizes, throughputs

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Process some CSV files.')
    parser.add_argument('result_path', type=str, help='The path of experiment result traces')
    parser.add_argument('test_name', type=str, help='The name of the test')
    args = parser.parse_args()
    
    version = "range_query"
    benchmark_log_dir = f'{args.result_path}/{args.test_name}'
    output_dir = f"benchmark_results_{version}"
    summary_file = f'{args.result_path}/{args.test_name}_summary.csv'
    summary_plot_file = f'{args.result_path}/{args.test_name}_summary.png'
    
    log_files = [log_file for log_file in os.listdir(benchmark_log_dir) if log_file.endswith('.log')]
    result_dict = {"entry_count": [], "range_size": [], "latency": [], "throughput": [], "timpestamp": []}
    for log_file in log_files:
        type, entry_count, range_sizes, latencies, throughputs, timestamp = parse_log(os.path.join(benchmark_log_dir, log_file))
        result_length = len(latencies)
        if type == "range-get" or type == "range-query":
            result_dict["entry_count"] += [entry_count] * result_length
            result_dict["range_size"] += range_sizes
            result_dict["latency"] += latencies
            result_dict["throughput"] += throughputs
            result_dict["timpestamp"] += [timestamp] * result_length

            
    result_df = pd.DataFrame(result_dict)
    result_df.sort_values(by=['entry_count', 'range_size'], ascending=[True,True], inplace=True, na_position='last')    
    summary_df = result_df.groupby(['entry_count', 'range_size'])[['latency','throughput']].mean().reset_index()
    summary_df.to_csv(summary_file, index=False)
    