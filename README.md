<a name="readme-top"></a>

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]

<!-- PROJECT LOGO -->
<br />
<div align="center">
<h3 align="center">Datadog Import</h3>

  <p align="center">
    This is a command line tool to convert other resources to Datadog
    <br />
    <a href="https://github.com/abruneau/datadog_import/issues">Report Bug</a>
    ·
    <a href="https://github.com/abruneau/datadog_import/issues">Request Feature</a>
  </p>
</div>

<!-- ABOUT THE PROJECT -->

## About The Project

### Built With

- [Go](https://go.dev/)
- [goreleaser](https://goreleaser.com/)
- [cobra](https://github.com/spf13/cobra)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->

## Getting Started

### Brew

```sh
brew install --cask abruneau/homebrew-tap/datadog_import
```

### From Binary

Get the [latest release](https://github.com/abruneau/datadog_import/releases) for your platform

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Build from source

1. Clone the repository
   ```sh
   git clone git@github.com:abruneau/datadog_import.git
   ```
2. Build
   ```sh
   go build -o datadog_import
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Usage

- Create a config file in `yaml`. You can look at [config.yaml.example](./config.yaml.example) for guidance.
- run the cli

### Root command

```
Usage:
  datadog_import [flags]
  datadog_import [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  grafana     Convert Grafana queries to Datadog queries
  help        Help about any command

Flags:
      --config string   config file (default is ./config.yaml)
  -h, --help            help for datadog_import
      --log string      log level (default is info)) (default "info")
  -v, --version         Print the version and exit
```

### Grafana command

```
Convert queries from various Grafana datasources (Azure Monitor, CloudWatch, Stackdriver, Prometheus, Loki) to Datadog queries.

Usage:
  datadog_import grafana [flags]

Flags:
  -h, --help                help for grafana
  -q, --query stringArray   Queries to convert (can be specified multiple times)
  -t, --type string         Query type (grafana-azure-monitor-datasource, cloudwatch, stackdriver, prometheus, loki)

Global Flags:
      --config string   config file (default is ./config.yaml)
      --log string      log level (default is info)) (default "info")
  -v, --version         Print the version and exit
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Compatibility

| Source    | Doc                       |
| --------- | ------------------------- |
| Grafana   | [doc](./doc/grafana.md)   |
| Dynatrace | [doc](./doc/dynatrace.md) |

<!-- ROADMAP -->

## Roadmap

See the [open issues](https://github.com/abruneau/datadog_import/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTRIBUTING -->

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- LICENSE -->

## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTACT -->

## Contact

Your Name - antonin.bruneau@gmail.com

Project Link: [https://github.com/abruneau/datadog_import](https://github.com/abruneau/datadog_import)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[contributors-shield]: https://img.shields.io/github/contributors/abruneau/datadog_import.svg?style=for-the-badge
[contributors-url]: https://github.com/abruneau/datadog_import/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/abruneau/datadog_import.svg?style=for-the-badge
[forks-url]: https://github.com/abruneau/datadog_import/network/members
[stars-shield]: https://img.shields.io/github/stars/abruneau/datadog_import.svg?style=for-the-badge
[stars-url]: https://github.com/abruneau/datadog_import/stargazers
[issues-shield]: https://img.shields.io/github/issues/abruneau/datadog_import.svg?style=for-the-badge
[issues-url]: https://github.com/abruneau/datadog_import/issues
[license-shield]: https://img.shields.io/github/license/abruneau/datadog_import.svg?style=for-the-badge
[license-url]: https://github.com/abruneau/datadog_import/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/antoninbruneau
