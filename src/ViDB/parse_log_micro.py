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
    batch_size = int(meta_info[2].strip())
    value_size = int(meta_info[3].strip())
    timestamp = meta_info[4].strip() + meta_info[5].strip()
    # print(f"type: {type}, entry_count: {entry_count}, batch_size: {batch_size}, value_size: {value_size}, timestamp: {timestamp}")
    if type == "micro":
        r_latency, w_latency, r_tps, w_tps = parse_micro_latencies(log_file, batch_size, entry_count)
        r_latency, w_latency, r_tps, w_tps = numpy.mean(r_latency), numpy.mean(w_latency), numpy.mean(r_tps), numpy.mean(w_tps)
        return type, entry_count, batch_size, value_size, timestamp, r_latency, w_latency, r_tps, w_tps


def parse_micro_latencies(log_file, batch_size, entry_count):
    r_latencies = []
    w_latencies = []
    r_tps = []
    w_tps = []
    with open(log_file, 'r') as f:
        for line in f:
            # print(line)
            match_w_batch = re.search(r'Execute Write (\d+) done.', line)
            match_r_batch = re.search(r'Execute Read (\d+) done.', line)
            if match_w_batch:
                assert int(match_w_batch.group(1)) == batch_size or int(match_w_batch.group(1)) == entry_count, "batch size not match"
                match_lat = re.search(r'Lantency: (\d+) ns', line)
                if match_lat:
                    w_latencies.append(int(match_lat.group(1)))
                match_tps = re.search(r'tps: ([\d.]+)', line)
                if match_tps:
                    w_tps.append(float(match_tps.group(1)))
            elif match_r_batch:
                assert int(match_r_batch.group(1)) == batch_size, "batch size not match"
                match_lat = re.search(r'Lantency: (\d+) ns', line)
                if match_lat:
                    r_latencies.append(int(match_lat.group(1)))
                match_qps = re.search(r'qps: ([\d.]+)', line)
                if match_qps:
                    r_tps.append(float(match_qps.group(1)))
            else:
                continue
    return r_latencies, w_latencies, r_tps, w_tps

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Process some CSV files.')
    parser.add_argument('result_path', type=str, help='The path of experiment result traces')
    parser.add_argument('test_name', type=str, help='The name of the test')
    args = parser.parse_args()
    
    benchmark_log_dir = f'{args.result_path}/{args.test_name}'
    summary_file = f'{args.result_path}/{args.test_name}_summary.csv'
    
    # list all files
    log_files = os.listdir(benchmark_log_dir)
    log_files = [os.path.join(benchmark_log_dir, log_file) for log_file in log_files if log_file.endswith('.log')]
    micro_result_dict = {"entry_count": [], "batch_size": [], "value_size": [], "read_latency": [], "write_latency": [], "read_throughput": [], "write_throughput": [], "timpestamp": []}
    for log_file in log_files:
        type, entry_count, batch_size, value_size, timestamp, r_latency, w_latency, r_tps, w_tps = parse_log(log_file)
        if type == "micro":
            micro_result_dict["entry_count"].append(entry_count)
            micro_result_dict["batch_size"].append(batch_size)
            micro_result_dict["value_size"].append(value_size)
            micro_result_dict["read_latency"].append(r_latency/1e9)
            micro_result_dict["write_latency"].append(w_latency/1e9)
            micro_result_dict["read_throughput"].append(r_tps)
            micro_result_dict["write_throughput"].append(w_tps)
            micro_result_dict["timpestamp"].append(timestamp)
    micro_result_df = pd.DataFrame(micro_result_dict)
    print(micro_result_df)
    micro_result_df.to_csv(f"{summary_file}", index=False)
