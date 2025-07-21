#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import re
import os
import argparse
import numpy
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
    batch_size = int(meta_info[2].strip())
    # value_size = int(meta_info[3].strip())
    timestamp = meta_info[-2].strip() + meta_info[-1].strip()
    print(f"type: {type}, entry_count: {entry_count}, batch_size: {batch_size}, timestamp: {timestamp}")
    if type == "dataLineage":
        ver_cnt, ver_ago, ver_start, ver_end, latencies_ns, latencies_s, throughputs = parse_lineage_latencies(log_file)
        return type, entry_count, ver_cnt, ver_ago, ver_start, ver_end, latencies_ns, latencies_s, throughputs, timestamp
    else:
        return "None", 0, [], [], [], [], [], [], [], timestamp
    
def parse_lineage_latencies(log_file):
 
    ver_cnt = []
    ver_ago = []
    ver_start = []
    ver_end = []
    latencies_ns = []
    latencies_s = []
    throughputs = []
    with open(log_file, 'r') as f:
        for line in f:
            pattern = r"query\s*(\d+)\s*versions\s*starting\s*from\s*(\d+)\s*versions\s*ago\s*\[(\d+)\s*,\s*(\d+)\),\s*Latency:\s*(\d+)\s*ns\s*\((\d+.\d+)\s*s\)\s*TPS:\s*(\d+.\d+)\(n/s\)"
            match = re.search(pattern, line)
            if match:
                ver_cnt.append(int(match.group(1)))
                ver_ago.append(int(match.group(2)))
                ver_start.append(int(match.group(3)))
                ver_end.append(int(match.group(4)))                
                latencies_ns.append(float(match.group(5)))
                latencies_s.append(float(match.group(6)))
                throughputs.append(float(match.group(7)))
            else:
                pass
    return ver_cnt, ver_ago, ver_start, ver_end, latencies_ns, latencies_s, throughputs

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Process some CSV files.')
    parser.add_argument('result_path', type=str, help='The path of experiment result traces')
    parser.add_argument('test_name', type=str, help='The name of the test')
    args = parser.parse_args()
    
    benchmark_log_dir = f'{args.result_path}/{args.test_name}'
    summary_file = f'{args.result_path}/{args.test_name}_summary.csv'
    summary_plot_file = f'{args.result_path}/{args.test_name}_summary.png'
    
    # 罗列文件夹内所有文件
    log_files = [log_file for log_file in os.listdir(benchmark_log_dir) if log_file.endswith('.log')]
    result_dict = {"entry_count":[], "ver_cnt": [], "ver_ago": [], "ver_start": [], "ver_end": [], "latencies_ns": [], "latencies_s": [], "throughputs": [], "timpestamp": []}
    for log_file in log_files:
        type, entry_count, ver_cnt, ver_ago, ver_start, ver_end, latencies_ns, latencies_s, throughputs, timestamp = parse_log(os.path.join(benchmark_log_dir, log_file))
        result_length = len(ver_cnt)
        if type == "dataLineage":
            result_dict["entry_count"] += [entry_count] * result_length
            result_dict["ver_cnt"] += ver_cnt
            result_dict["ver_ago"] += ver_ago
            result_dict["ver_start"] += ver_start
            result_dict["ver_end"] += ver_end
            result_dict["latencies_ns"] += latencies_ns
            result_dict["latencies_s"] += latencies_s
            result_dict["throughputs"] += throughputs
            result_dict["timpestamp"] += [timestamp] * result_length
            
    result_df = pd.DataFrame(result_dict)
    result_df.sort_values(by=['entry_count', 'ver_cnt'], ascending=[True,True], inplace=True, na_position='last')
    print(result_df)
    summary_df = result_df.groupby(['entry_count', 'ver_cnt', 'ver_ago'])[['latencies_s','throughputs']].mean().reset_index()
    summary_df.to_csv(summary_file, index=False)
    