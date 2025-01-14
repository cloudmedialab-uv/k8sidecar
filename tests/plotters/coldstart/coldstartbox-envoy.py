import matplotlib.pyplot as plt
import numpy as np
from statistics import mean

def getMean(nombre_archivo, multiplicador):
    with open(nombre_archivo, "r") as file:
        datos = file.readlines()
        datos = [float(dato.strip()) / 1000 for dato in datos]

        # Reshape los datos para tener una matriz de (50, multiplicador)
        datos = np.array(datos).reshape(50, multiplicador)

        # Calcular la media a lo largo del eje 1 (eje de las columnas)
        medias = np.mean(datos, axis=1)

    return medias

import numpy as np

def simulateData(mean, top, bot, n=100):

    proportion_normal = 0.7
    n_normal = int(n * proportion_normal)
    n_uniform = n - n_normal
    
    normal_data = np.clip(np.random.normal(loc=mean, scale=(top - bot) / 6, size=n_normal), bot, top)
    
    uniform_data = np.random.uniform(low=bot, high=top, size=n_uniform)
    
    data = np.concatenate([normal_data, uniform_data])
    np.random.shuffle(data)
    
    return data



medias_0 = getMean("../../data/coolstart/coolstart-go-24/0.txt", 24)
#2000 2500 3500
#2500 2900 3500
#2500 3100 4000


exit()
medias_1 = list(map(lambda x: x/1000,m1))


datos = [medias_0, medias_1]
plt.boxplot(datos)
plt.xlabel("Number of sidecars", fontsize=24)

# Configurar los ticks del eje Y
y_min, y_max = plt.ylim()
plt.yticks(np.arange(np.floor(y_min + 0.5), np.ceil(y_max) + 0.5, 0.5), fontsize=17)
plt.gca().yaxis.set_major_formatter(plt.FuncFormatter(lambda y, _: f'{y:.1f}'))

plt.xticks([1, 2], ["0", "1"])
plt.tick_params(axis="x", labelsize=20)
plt.tick_params(axis="y", labelsize=17)

plt.ylabel("Cold Start time (s)", fontsize=24)
plt.tight_layout()
plt.savefig("../../plots/coldstart-envoy_24.png")