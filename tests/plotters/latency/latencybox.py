import matplotlib.pyplot as plt
import json
from statistics import mean

def getLatency(name):
    with open(name, "r") as file:
        times=[]
        for line in file:
            data = json.loads(line)
            latency= data["tsFinalTime"]-data["tsReceivedTime"]
            if latency < 40:
                times.append(latency)

    return times


dir = "../../data/latency-go/"

# Obtiene la media para cada archivo
datos_0 = getLatency(dir + "0.json")
datos_1 = getLatency(dir + "1.json")
datos_2 = getLatency(dir + "2.json")
datos_5 = getLatency(dir + "5.json")
datos_10 = getLatency(dir + "10.json")

datos = [datos_0, datos_1, datos_2, datos_5, datos_10]

print("media de datos_0: ", mean(datos_0))
print("media de datos_1: ", mean(datos_1))
print("media de datos_2: ", mean(datos_2))
print("media de datos_5: ", mean(datos_5))
print("media de datos_10: ", mean(datos_10))

print("diferencia entre 0 y 10", mean(datos_10)-mean(datos_0))


plt.boxplot(datos)
plt.xticks(
    [1, 2, 3, 4, 5], ["0 Sidecar", "1 Sidecar", "2 Sidecar", "5 Sidecar", "10 Sidecar"]
)
plt.ylabel("Latency time (ms)")
plt.show()