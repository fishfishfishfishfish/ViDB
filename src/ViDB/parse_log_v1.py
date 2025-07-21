import re
import os

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
    version = "V3"
    benchmark_log_dir = f"benchmark_results_{version}/"
    summary_file = f"summary.csv"
    # 罗列文件夹内所有文件
    log_files = os.listdir(benchmark_log_dir)
    log_files = [os.path.join(benchmark_log_dir, log_file) for log_file in log_files if log_file.endswith('.log')]
    rw_result_dict = {"entry_count": [], "batch_size": [], "value_size": [], "operation": [], "latency": [], "timpestamp": []}
    micro_result_dict = {"entry_count": [], "batch_size": [], "value_size": [], "read_latency": [], "write_latency": [], "timpestamp": []}
    for log_file in log_files:
        type, entry_count, batch_size, value_size, timestamp, latency = parse_log(log_file)
        if type == "micro":
            micro_result_dict["entry_count"].append(entry_count)
            micro_result_dict["batch_size"].append(batch_size)
            micro_result_dict["value_size"].append(value_size)
            micro_result_dict["read_latency"].append(latency[0])
            micro_result_dict["write_latency"].append(latency[1])
            micro_result_dict["timpestamp"].append(timestamp)
        else:
            rw_result_dict["entry_count"].append(entry_count)
            rw_result_dict["batch_size"].append(batch_size)
            rw_result_dict["value_size"].append(value_size)
            rw_result_dict["operation"].append(type)
            rw_result_dict["latency"].append(latency)
            rw_result_dict["timpestamp"].append(timestamp)
    rw_result_df = pd.DataFrame(rw_result_dict)
    micro_result_df = pd.DataFrame(micro_result_dict)
    print(rw_result_df)
    print(micro_result_df)
    micro_result_df.to_csv(f"summary_{version}.csv", index=False)
    rw_result_df.to_csv(f"restart-summary_{version}.csv", index=False)
    
    plt.figure(figsize=(4, 3))
    plt.plot(micro_result_df["batch_size"], micro_result_df["read_latency"], label="no restart")
    plt.plot(rw_result_df["batch_size"], rw_result_df["latency"], label="restart db")
    plt.ylabel("latency (ns)")
    plt.xlabel("value size")
    plt.title("Read latency")
    plt.legend()
    
    plt.savefig(f"read_latency_{version}.png")