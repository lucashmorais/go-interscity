import sys
import math
import numpy as np
import random
from matplotlib import pyplot as plt

raw_data = []

for line in sys.stdin:
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

plt.savefig("test.svg", format="svg")
plt.savefig("test.png", format="png", dpi=200)