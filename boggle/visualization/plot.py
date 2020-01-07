#!/usr/bin/python3

import pandas as pd
import matplotlib as mpl
import matplotlib.pyplot as plt

df = pd.read_csv('boggle.csv', header=None,
    names=['time','score'], usecols=[1,2])
x = df['time'] / 1000 # seconds
y = df['score'] # points

plt.figure(0)
plt.semilogx(x, y, color="#c05131", lw=2)
plt.hlines(4540, xmin=1e-1, xmax=1e17, color="#6c6f70", lw=1, linestyle='dashed')
plt.plot([x.max(), 1e17], [y.max(), 4540], color="#ef8200", lw=2, linestyle='dotted')
plt.ylabel("best score so far")
plt.xlabel("simulation time (s)")
plt.grid(b=True, which='major', axis='x', color="#6c6f70", lw=.5, linestyle='dotted')

boxstyle = {'fc':'white', 'lw':0, 'boxstyle':'round4'}
plt.text(.5, 2000, 'my simulations', color="#c05131", bbox=boxstyle)
plt.text(.5, 4300, 'best score from Sedgewick & Wayne students', color="#6c6f70", bbox=boxstyle)
plt.text(1e10, 2500, 'performance required\nto beat S&W students\nbefore sun goes nova', color="#ef8200", bbox=boxstyle)

plt.title("Plot to convince my wife I need to upgrade my computer")

plt.savefig('score.png')
