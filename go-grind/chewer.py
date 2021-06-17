#!/usr/bin/env python3
import sys
import os
import signal
import math
import numpy as np
import random
import time
import argparse
from subprocess import Popen
from matplotlib import pyplot as plt

def process_log(num_clients: int, num_requests: int):
	raw_data = []
	
	base_log_path = "/home/lucas/Repos/go-interscity/resource-adaptor/"
	
	log_path = base_log_path + log_name(num_clients, num_requests)
	file = open(log_path)

	for line in file.readlines():
		if line.startswith('['):
			number_string = line.split(' ')[1]
			number = float(number_string) / 1000
			raw_data.append(number)

	bins = np.linspace(math.ceil(min(raw_data)), 
			   math.floor(max(raw_data)),
			   1000)

	plt.xlim([min(raw_data)-5, max(raw_data)+5])
	plt.hist(raw_data, bins=bins, alpha=0.5)
	plt.title('Server latency')
	plt.xlabel('Latency (milliseconds)')
	plt.ylabel('Request count')

	# print("Average latency (milliseconds): %d" % np.average(raw_data))
	# print("Latency standard deviation (milliseconds): %d" % np.std(raw_data))
	
	base_figure_path = "/home/lucas/Repos/go-interscity/go-grind/output/"

	plt.savefig(f"{base_figure_path}{num_clients}_{num_requests}.svg", format="svg")
	plt.savefig(f"{base_figure_path}{num_clients}_{num_requests}.png", format="png", dpi=200)
	plt.close()
	
def log_name(num_clients, num_requests_per_client):
	return f"{num_clients}_{num_requests_per_client}.log"
    
def get_server_command(num_clients: int, num_requests: int):
	base = "cd /home/lucas/Repos/go-interscity/resource-adaptor/ && go run server.go > " 
	command = base + f"{log_name(num_clients, num_requests)} 2>&1"
	return command
	
def get_grinder_command(num_clients: int, num_requests: int):
	command = f"cd /home/lucas/Repos/go-interscity/go-grind && sleep 2 && go run grinder.go {num_clients} {num_requests}"
	return command
	
def core_spawn_test(num_clients: int, num_requests: int):
	commands = []

	grinder_command = get_grinder_command(num_clients, num_requests)
	server_command = get_server_command(num_clients, num_requests)

	commands.append(grinder_command)
	commands.append(server_command)
	
	procs = [ Popen(args=i, shell=True, preexec_fn=os.setsid) for i in commands ]

	# https://stackoverflow.com/questions/4789837/how-to-terminate-a-python-subprocess-launched-with-shell-true
	print(procs[0].args)
	procs[0].wait()
	procs[1].kill()
	print()

	os.killpg(os.getpgid(procs[1].pid), signal.SIGTERM)
	
def get_split_set(max_value: int, num_values: int):
	step = float(max_value) / float(num_values)
	set = [ round(step * i) for i in range(1, num_values) ]
	set.append(max_value)
	return set
	
def get_set_of_num_clients(max_num_clients: int, num_items: int):
	return get_split_set(max_num_clients, num_items)
	
def spawn_test(args):
	max_num_clients = args.num_clients
	num_requests = args.requests_per_client
	num_tests = 10
	
	set_of_num_clients = get_set_of_num_clients(max_num_clients, num_tests)
	for num_clients in set_of_num_clients:
		core_spawn_test(num_clients, num_requests)
		process_log(num_clients, num_requests)

# https://zetcode.com/python/argparse/
argument_parser = argparse.ArgumentParser()

argument_parser.add_argument('--num-clients', type=int, required=True)
argument_parser.add_argument('--requests-per-client', type=int, required=True)
argument_parser.add_argument('--driver', dest='driver', choices=['requests', 'request'], help="Defines which test driver to use", default='requests')
argument_parser.add_argument('--uuid', type=str, required=False)

args = argument_parser.parse_args()

spawn_test(args)