#!/usr/bin/env python3
import sys
import os
import signal
import math
import numpy as np
import random
import time
import argparse
import requests
from subprocess import Popen
from matplotlib import pyplot as plt

SERVER_BASE_ROUTE = 'http://localhost'
SERVER_TEST_ROUTE = '/'
SERVER_PORT = 8888

class LatencyInfo:
	def __init__(self, minimum, q10, med, average, q90, q95, std):
		self.minimum = minimum
		self.q10 = q10
		self.med = med
		self.average = average
		self.q90 = q90
		self.q95 = q95
		self.std = std
	def __str__(self):
		return f"{self.minimum}, {self.q10}, {self.med}, {self.average}, {self.q90}, {self.q95}, {self.std}"
	
def get_test_route():
	return f"{SERVER_BASE_ROUTE}:{SERVER_PORT}{SERVER_TEST_ROUTE}"
	
def test_server_is_up():
	try:
		testRoute = get_test_route()
		r = requests.get(testRoute)
		return r.status_code == 200
	except:
		return False

def wait_until_server_is_up():
	while (not test_server_is_up()):
		time.sleep(0.3)

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
	
	average = np.average(raw_data)
	std = np.std(raw_data)
	median = np.median(raw_data)
	quantiles = np.quantile(raw_data, [0.1, 0.9, 0.95])
	minimum= np.min(raw_data)

	latencyInfo = LatencyInfo(minimum, quantiles[0], median, average, quantiles[1], quantiles[2], std)
	print(latencyInfo)

	return latencyInfo
	
def log_name(num_clients, num_requests_per_client):
	return f"{num_clients}_{num_requests_per_client}.log"
    
def get_server_command(num_clients: int, num_requests: int):
	base = "cd /home/lucas/Repos/go-interscity/resource-adaptor/ && go run server.go > " 
	command = base + f"{log_name(num_clients, num_requests)} 2>&1"
	return command
	
def get_grinder_command(num_clients: int, num_requests: int):
	command = f"cd /home/lucas/Repos/go-interscity/go-grind && go run grinder.go {num_clients} {num_requests}"
	return command

def get_freq_set_command(target_frequency: int):
	command = f"/home/lucas/Repos/go-interscity/go-grind/max_perf_custom_freq.sh {target_frequency}"
	return command
	
def set_frequency(target_frequency: int):
	print(f"Setting frequency to {target_frequency} MHz")
	command = get_freq_set_command(target_frequency)
	proc = Popen(args=command, shell=True, preexec_fn=os.setsid)
	proc.wait()
	print(f"Frequency was successfully set to {target_frequency} MHz")
	
def core_spawn_test(num_clients: int, num_requests: int):
	commands = []

	grinder_command = get_grinder_command(num_clients, num_requests)
	server_command = get_server_command(num_clients, num_requests)

	server_proc = Popen(args=server_command, shell=True, preexec_fn=os.setsid)
	wait_until_server_is_up()
	grinder_proc = Popen(args=grinder_command, shell=True, preexec_fn=os.setsid)

	commands.append(grinder_command)
	commands.append(server_command)
	
	procs = [ grinder_proc, server_proc ]

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

def plotLatencyInfo(set_of_num_clients, data):
	averages = [ d.average for d in data ]
	errors = [ d.std for d in data ]
	mins = [ d.minimum for d in data ]
	p10 = [ d.q10 for d in data ]
	meds = [ d.med - d.q10 for d in data ]
	true_meds = [ d.med for d in data ]
	p90 = [ d.q90 - d.med for d in data ]
	true_p90 = [ d.q90 for d in data ]
	p95 = [ d.q95 for d in data ]
	print(set_of_num_clients)
	print(p10)
	print(meds)

	# Data for plotting
	# t = np.arange(0.0, 2.0, 0.01)
	# s = 1 + np.sin(2 * np.pi * t)

	t = set_of_num_clients
	s = averages

	fig, ax = plt.subplots()
	# ax.plot(t, s)
	# ax.errorbar(t, s, errors)
	# ax.errorbar(t, s, errors, marker='^', capsize=3)
	# ax.errorbar(t, s, marker='^', capsize=3)
	ax.errorbar(t, s, capsize=20, label="avg")

	ax.grid()
	ax.set_axisbelow(True)

	_, caps0, _ = ax.errorbar(t, p10, [ d.q10 - d.minimum for d in data ], linestyle='', uplims=True, capsize=3, ecolor='#457B9D')
	_, caps1, _ = ax.errorbar(t, true_p90, [ d.q95 - d.q90 for d in data ], linestyle='', lolims=True, capsize=3, ecolor='#E63946')

	# ax.errorbar(t, p10, [[ d.q10 - d.minimum for d in data ], [ 0 for d in data ]], linestyle='', capsize=3, ecolor='blue')
	# ax.errorbar(t, true_p90, [[0 for d in data], [ d.q95 - d.q90 for d in data ]], linestyle='', capsize=3, ecolor='green')

	# https://stackoverflow.com/questions/45752981/removing-the-bottom-error-caps-only-on-matplotlib
	ax.bar(t, meds, width=100, bottom=p10, color="#5590B4",  edgecolor='#457B9D', label="min/p10/med")
	ax.bar(t, p90, width=100, bottom=true_meds, color="#F07F89",  edgecolor='#E63946', label="med/p90/p95")
	
	caps0[0].set_marker('_')
	caps0[0].set_markersize(10)
	caps1[0].set_marker('_')
	caps1[0].set_markersize(10)

	ax.set(xlabel='Number of concurrent clients', ylabel='Latency (milliseconds)',
	       title='Average latency')
	# plt.yticks(np.arange(0, 1000, step=0.1))  # Set label locations.
	
	# Information on how to build custom scales
	# https://stackoverflow.com/questions/31168051/creating-probability-frequency-axis-grid-irregularly-spaced-with-matplotlib/31170170#31170170
	plt.legend(loc='upper left')

	base_figure_path = "/home/lucas/Repos/go-interscity/go-grind/output/"

	fig.savefig(f"{base_figure_path}degradation_plot_up_to_{set_of_num_clients[-1]}_clients.png", dpi=200)

	plt.yscale('log')
	fig.savefig(f"{base_figure_path}degradation_plot_up_to_{set_of_num_clients[-1]}_clients_log.png", dpi=200)

	plt.show()
	plt.close()
	
def spawn_test(args):
	max_num_clients = args.num_clients
	num_requests = args.requests_per_client
	skip_measurements = args.skip_measurements
	num_tests = args.num_steps
	min_freq = args.min_cpu_freq
	max_freq = args.max_cpu_freq
	
	if (min_freq == max_freq):
		set_frequency(min_freq)
	
	set_of_num_clients = get_set_of_num_clients(max_num_clients, num_tests)
	latencyInfo = []
	for num_clients in set_of_num_clients:
		if (not skip_measurements):
			core_spawn_test(num_clients, num_requests)
		latencyInfo.append(process_log(num_clients, num_requests))
		
	plotLatencyInfo(set_of_num_clients, latencyInfo)

# https://zetcode.com/python/argparse/
argument_parser = argparse.ArgumentParser()

argument_parser.add_argument('--num-clients', type=int, required=True)
argument_parser.add_argument('--num-steps', type=int, default=10, required=False)
argument_parser.add_argument('--requests-per-client', type=int, required=True)
argument_parser.add_argument('--driver', dest='driver', choices=['requests', 'request'], help="Defines which test driver to use", default='requests')
argument_parser.add_argument('--uuid', type=str, required=False)
argument_parser.add_argument('--min-cpu-freq', type=int, required=False, help="Defines minimum processor frequency to be tested", default=3400)
argument_parser.add_argument('--max-cpu-freq', type=int, required=False, help="Defines maximum processor frequency to be tested", default=3400)

# https://docs.python.org/3/howto/argparse.html#introducing-optional-arguments
# The "store_true" action makes argparser automatically assign "True" to the related variable anytime
# the optional argument is found, avoiding us to require the user to pass any value after setting the flag 
argument_parser.add_argument('--skip-measurements', default=False, required=False, action="store_true")

args = argument_parser.parse_args()
print(args)
# print(args.skip_measurements)

spawn_test(args)