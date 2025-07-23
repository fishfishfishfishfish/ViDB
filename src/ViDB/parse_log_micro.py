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
    if type == "random-get":
        latency = numpy.mean(parse_get_latencies(log_file, batch_size))
        return type, entry_count, batch_size, value_size, timestamp, latency
    elif type == "random-write":
        latency = numpy.mean(parse_write_latencies(log_file, batch_size, entry_count))
        return type, entry_count, batch_size, value_size, timestamp, latency
    elif type == "micro":
        r_latency, w_latency = parse_micro_latencies(log_file, batch_size, entry_count)
        r_latency, w_latency = numpy.mean(r_latency), numpy.mean(w_latency)
        return type, entry_count, batch_size, value_size, timestamp, (r_latency, w_latency)
    
def parse_get_latencies(log_file, batch_size):
    latencies = []
    with open(log_file, 'r') as f:
        for line in f:
            match_batch = re.search(r'Execute Read (\d+) done.', line)
            if match_batch:
                assert int(match_batch.group(1)) == batch_size, "batch size not match"
                match_lat = re.search(r'Lantency: (\d+) ns', line)
                if match_lat:
                    latencies.append(int(match_lat.group(1)))
            else:
                continue
    return latencies

def parse_write_latencies(log_file, batch_size, entry_count):
    latencies = []
    with open(log_file, 'r') as f:
        for line in f:
            match_batch = re.search(r'Execute Write (\d+) done.', line)
            if match_batch:
                assert int(match_batch.group(1)) == batch_size or int(match_batch.group(1)) == entry_count, "batch size not match"
                match_lat = re.search(r'Lantency: (\d+) ns', line)
                if match_lat:
                    latencies.append(int(match_lat.group(1)))
            else:
                continue
    return latencies

def parse_micro_latencies(log_file, batch_size, entry_count):
    r_latencies = []
    w_latencies = []
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
            elif match_r_batch:
                assert int(match_r_batch.group(1)) == batch_size, "batch size not match"
                match_lat = re.search(r'Lantency: (\d+) ns', line)
                if match_lat:
                    r_latencies.append(int(match_lat.group(1)))
            else:
                continue
    return r_latencies, w_latencies

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Process some CSV files.')
    parser.add_argument('result_path', type=str, help='The path of experiment result traces')
    parser.add_argument('test_name', type=str, help='The name of the test')
    args = parser.parse_args()
    
    benchmark_log_dir = f'{args.result_path}/{args.test_name}'
    summary_file = f'{args.result_path}/{args.test_name}_summary.csv'
    summary_plot_file = f'{args.result_path}/{args.test_name}_summary.png'
    
    # list all files
    log_files = os.listdir(benchmark_log_dir)
    log_files = [os.path.join(benchmark_log_dir, log_file) for log_file in log_files if log_file.endswith('.log')]
    micro_result_dict = {"entry_count": [], "batch_size": [], "value_size": [], "get_latency": [], "put_latency": [], "get_throughput": [], "put_throughput": [], "timpestamp": []}
    for log_file in log_files:
        type, entry_count, batch_size, value_size, timestamp, latency = parse_log(log_file)
        if type == "micro":
            micro_result_dict["entry_count"].append(entry_count)
            micro_result_dict["batch_size"].append(batch_size)
            micro_result_dict["value_size"].append(value_size)
            micro_result_dict["get_latency"].append(latency[0]*1e-9)
            micro_result_dict["put_latency"].append(latency[1]*1e-9)
            micro_result_dict["get_throughput"].append(batch_size/(latency[0]*1e-9))
            micro_result_dict["put_throughput"].append(batch_size/(latency[1]*1e-9))
            micro_result_dict["timpestamp"].append(timestamp)
    micro_result_df = pd.DataFrame(micro_result_dict)
    micro_result_df.sort_values(by=['entry_count', 'batch_size', 'value_size'], ascending=[True,True,True], inplace=True, na_position='last')    

    print(micro_result_df)
    micro_result_df.to_csv(summary_file, index=False)
    
    plt.figure(figsize=(4, 3))
    plt.plot(micro_result_df["batch_size"], micro_result_df["get_latency"], label="no restart")
    plt.ylabel("latency (s)")
    plt.xlabel("value size")
    plt.title("Read latency")
    plt.legend()
    
    plt.savefig(summary_plot_file)