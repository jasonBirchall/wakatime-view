import typer
import base64
import requests
from rich import print


app = typer.Typer()
wakatime_api = "https://wakatime.com/api/v1/"


@app.command()
def today():
    api_key = typer.prompt("Enter your API key: ")

    endpoint = "users/current/status_bar/today"
    base64_bytes = base64.b64encode(api_key.encode("ascii"))
    base64_api_key = base64_bytes.decode("ascii")
    authorization = f"Basic {base64_api_key}"
    headers = {"Authorization": authorization}
    req = requests.get(wakatime_api + endpoint, headers=headers)

    if req.status_code == 200:
        response = req.json()
        print(f"{response['data']['grand_total']['text']}")
    else:
        print("Error", req)


@app.command()
def goodbye(name: str, formal: bool = False):
    if formal:
        typer.echo(f"Goodbye Ms. {name}. Have a good day.")
    else:
        typer.echo(f"Bye {name}!")


if __name__ == "__main__":
    app()
