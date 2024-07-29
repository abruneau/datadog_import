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
<h3 align="center">Dynatrace synthetics to Datadog</h3>

  <p align="center">
    This is a command line tool to convert Dynatrace synthetics to Datadog synthetics
    <br />
    <a href="https://github.com/abruneau/dynatrace_to_datadog/issues">Report Bug</a>
    ·
    <a href="https://github.com/abruneau/dynatrace_to_datadog/issues">Request Feature</a>
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
brew install abruneau/homebrew-tap/dynatrace_to_datadog
```

### From Binary

Get the [latest release](https://github.com/abruneau/dynatrace_to_datadog/releases) for your platform

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Build from source

1. Clone the repository
   ```sh
   git clone git@github.com:abruneau/dynatrace_to_datadog.git
   ```
2. Build
   ```sh
   go build -o dynatrace_to_datadog
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Usage

- Create a config file in `yaml`. You can look at [config.yaml.example](./config.yaml.example) for guidance.
- run the cli

```sh
Usage:
  dynatrace_to_datadog [flags]

Flags:
      --config string   config file (default is ./config.yaml)
  -h, --help            help for dynatrace_to_datadog
      --log string      log level (default is info)) (default "info")
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Compatibility

| Source  | Support |
| ------- | ------- |
| browser | ✅      |
| api     | ✅      |

### API Tests Limitations

- Pre and post scripts don't have equivalent in Datadog. The details of those scripts will be reported in the message section of the test.
- Certificate Expiry Date Constraint is not supported in HTTP tests. It should be implemented in SSL tests

<!-- ROADMAP -->

## Roadmap

See the [open issues](https://github.com/abruneau/dynatrace_to_datadog/issues) for a full list of proposed features (and known issues).

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

Project Link: [https://github.com/abruneau/dynatrace_to_datadog](https://github.com/abruneau/dynatrace_to_datadog)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[contributors-shield]: https://img.shields.io/github/contributors/abruneau/dynatrace_to_datadog.svg?style=for-the-badge
[contributors-url]: https://github.com/abruneau/dynatrace_to_datadog/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/abruneau/dynatrace_to_datadog.svg?style=for-the-badge
[forks-url]: https://github.com/abruneau/dynatrace_to_datadog/network/members
[stars-shield]: https://img.shields.io/github/stars/abruneau/dynatrace_to_datadog.svg?style=for-the-badge
[stars-url]: https://github.com/abruneau/dynatrace_to_datadog/stargazers
[issues-shield]: https://img.shields.io/github/issues/abruneau/dynatrace_to_datadog.svg?style=for-the-badge
[issues-url]: https://github.com/abruneau/dynatrace_to_datadog/issues
[license-shield]: https://img.shields.io/github/license/abruneau/dynatrace_to_datadog.svg?style=for-the-badge
[license-url]: https://github.com/abruneau/dynatrace_to_datadog/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/antoninbruneau
