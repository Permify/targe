<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://github.com/Permify/kivo/raw/master/assets/images/logo-iam-copilot-dark.png">
    <img alt="Kivo logo" src="https://github.com/Permify/kivo/raw/master/assets/images/logo-iam-copilot-light.png" width="40%">
  </picture>
<h1 align="center">
   Kivo - Open Source IAM Copilot
</h1>
</div>
<p align="center">
   Kivo is an open-source CLI for managing IAM (Identity and Access Management) operations with AI assistance.
</p>

![Kivo Demo](assets/images/kivo.gif)

<p align="center">
    <a href="https://github.com/Permify/kivo" target="_blank"><img src="https://img.shields.io/github/go-mod/go-version/Permify/kivo?style=for-the-badge&logo=go" alt="Permify Go Version" /></a>&nbsp;
 <a href="https://goreportcard.com/report/github.com/Permify/kivo" target="_blank"><img src="https://goreportcard.com/badge/github.com/Permify/kivo?style=for-the-badge&logo=go" alt="Kivo Go Report Card" /></a>&nbsp;
    <a href="https://github.com/Permify/kivo" target="_blank"><img src="https://img.shields.io/github/license/Permify/kivo?style=for-the-badge" alt="Kivo Licence" /></a>&nbsp;
    <a href="https://discord.gg/n6KfzYxhPp" target="_blank"><img src="https://img.shields.io/discord/950799928047833088?style=for-the-badge&logo=discord&label=DISCORD" alt="Permify Discord Channel" /></a>&nbsp;
    <a href="https://github.com/Permify/kivo/releases" target="_blank"><img src="https://img.shields.io/github/v/release/permify/kivo?include_prereleases&style=for-the-badge" alt="Kivo Release" /></a>&nbsp;
    <a href="https://img.shields.io/github/commit-activity/m/Permify/kivo?style=for-the-badge" target="_blank"><img src="https://img.shields.io/github/commit-activity/m/Permify/kivo?style=for-the-badge" alt="Kivo Commit Activity" /></a>&nbsp;
    <a href="https://img.shields.io/github/actions/workflow/status/Permify/kivo/release.yml?style=for-the-badge" target="_blank"><img src="https://img.shields.io/github/actions/workflow/status/Permify/kivo/release.yml?style=for-the-badge" alt="GitHub Workflow Status" /></a>&nbsp;
</p>     

<p align="center">
   DevOps engineers use Kivo to configure how employees in their organization access infrastructure resources. Kivo simplifies and accelerates granting, revoking, and creating custom policies within the CLI — eliminating the need for tedious back-and-forth UI work.
</p>

### Installation Steps

1. **Install Kivo CLI:**
   ```shell
   brew tap permify/tap-kivo
   brew install kivo
   ```

2. **Set Up AWS Credentials:**

   Kivo requires AWS credentials to be configured in the file `~/.aws/credentials`. Follow these steps:

    - Create or open the `~/.aws/credentials` file using a text editor:
      ```shell
      nano ~/.aws/credentials
      ```

    - Add your AWS credentials in the following format:
      ```plaintext
      [default]
      aws_access_key_id = your_access_key
      aws_secret_access_key = your_secret_key
      ```

    - Save the file and exit (in nano, press `CTRL + O` to save, then `CTRL + X` to exit).

3. **Verify the Configuration:**

   Run the following command to confirm the credentials are set correctly:
   ```shell
   aws sts get-caller-identity
   ```
   This should return information about your AWS account. If it fails, double-check the credentials file for accuracy.

4. **Set the Default Region (Optional):**

   If your tool requires a specific AWS region, you can set it in the `~/.aws/config` file:
   ```shell
   nano ~/.aws/config
   ```
   Add:
   ```plaintext
   [default]
   region = us-east-1
   ```
   Replace `us-east-1` with your desired region.

## Communication Channels

If you like Permify, please consider giving us a :star:

<p align="left">
<a href="https://discord.gg/n6KfzYxhPp">
 <img height="70px" width="70px" alt="permify | Discord" src="https://user-images.githubusercontent.com/39353278/187209316-3d01a799-c51b-4eaa-8f52-168047078a14.png" />
</a>
<a href="https://twitter.com/GetPermify">
  <img height="70px" width="70px" alt="permify | Twitter" src="https://user-images.githubusercontent.com/39353278/187209323-23f14261-d406-420d-80eb-1aa707a71043.png"/>
</a>
<a href="https://www.linkedin.com/company/permifyco">
  <img height="70px" width="70px" alt="permify | Linkedin" src="https://user-images.githubusercontent.com/39353278/187209321-03293a24-6f63-4321-b362-b0fc89fdd879.png" />
</a>
</p>
