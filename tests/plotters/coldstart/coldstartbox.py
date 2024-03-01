import matplotlib.pyplot as plt
import numpy as np


def getMean(nombre_archivo, multiplicador):
    with open(nombre_archivo, "r") as file:
        datos = file.readlines()
        datos = [float(dato.strip()) / 1000 for dato in datos]

        # Reshape los datos para tener una matriz de (50, multiplicador)
        datos = np.array(datos).reshape(50, multiplicador)

        # Calcular la media a lo largo del eje 1 (eje de las columnas)
        medias = np.mean(datos, axis=1)

    return medias


dir = "../../data/coolstart/coolstart-java-72/"

# Obtiene la media para cada archivo
medias_0 = getMean(dir + "0.txt", 72)
medias_1 = getMean(dir + "1.txt", 72)
medias_2 = getMean(dir + "2.txt", 72)
medias_5 = getMean(dir + "5.txt", 72)
medias_10 = getMean(dir + "10.txt", 72)

datos = [medias_0, medias_1, medias_2, medias_5, medias_10]

plt.boxplot(datos)
plt.xlabel("Number of sidecars", fontsize=24)

plt.xticks(
    [1, 2, 3, 4, 5], ["0", "1", "2", "5", "10"]
)
plt.tick_params(axis="x",labelsize=20)
plt.ylabel("Cold Start time (s)", fontsize=24)
plt.show()
