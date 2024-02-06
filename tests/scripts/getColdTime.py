import json
import sys
import requests
import os


def main(data_file_url, output_file):
    # Realizar una petición GET al archivo de datos
    response = requests.get(data_file_url)

    # Verificar que la respuesta es exitosa
    if response.status_code != 200:
        print(
            f"Error al obtener el archivo de datos desde {data_file_url}. Código de estado: {response.status_code}"
        )
        return

    # Convertir la respuesta en líneas para iterar sobre ellas
    lines = response.text.splitlines()

    with open(output_file, "w") as out_f:
        for line in lines:
            try:
                # Convertir la línea a un objeto JSON
                data = json.loads(line)

                # Extraer tsReceivedTime y tsEvenGeneratedTime
                ts_received_time = data["tsReceivedTime"]
                ts_even_generated_time = data["tsEvenGeneratedTime"]

                # Calcular la resta y escribir en el archivo de salida
                resta = ts_received_time - ts_even_generated_time
                out_f.write(str(resta) + "\n")

            except Exception as e:
                print(f"Error procesando línea: {line}. Error: {e}")


if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Uso: script_name.py DATA_FILE_URL OUTPUT_FILE")
        sys.exit(1)

    # Obtener la variable de entorno UPLOAD_URL
    upload_url = os.environ.get("UPLOAD_SERVER_URL")

    if not upload_url:
        print("La variable de entorno UPLOAD_URL no está definida.")
        sys.exit(1)

    data_file_url = f"{upload_url}/{sys.argv[1]}"
    output_file = sys.argv[2]

    main(data_file_url, output_file)
