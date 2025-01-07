<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://github.com/Permify/targe/raw/master/assets/images/logo-targe-dark.png">
    <img alt="Targe logo" src="https://github.com/Permify/targe/raw/master/assets/images/logo-targe-light.png" width="40%">
  </picture>
<h1 align="center">
   Targe - Open Source IAM Copilot
</h1>
</div>
<p>
Targe is an open-source CLI for managing IAM (Identity and Access Management) operations with AI assistance.

DevOps engineers use Targe to configure how employees in their organization access infrastructure resources. Targe simplifies and accelerates granting and revoking access, while supporting custom policy creation — eliminating the need for tedious back-and-forth UI work.
</p>

<p align="center">
    <a href="https://github.com/Permify/targe" target="_blank"><img src="https://img.shields.io/github/go-mod/go-version/Permify/targe?style=for-the-badge&logo=go" alt="Permify Go Version" /></a>&nbsp;
 <a href="https://goreportcard.com/report/github.com/Permify/targe" target="_blank"><img src="https://goreportcard.com/badge/github.com/Permify/targe?style=for-the-badge&logo=go" alt="Targe Go Report Card" /></a>&nbsp;
    <a href="https://github.com/Permify/targe" target="_blank"><img src="https://img.shields.io/github/license/Permify/targe?style=for-the-badge" alt="Targe Licence" /></a>&nbsp;
    <a href="https://discord.gg/n6KfzYxhPp" target="_blank"><img src="https://img.shields.io/discord/950799928047833088?style=for-the-badge&logo=discord&label=DISCORD" alt="Permify Discord Channel" /></a>&nbsp;
    <a href="https://github.com/Permify/targe/releases" target="_blank"><img src="https://img.shields.io/github/v/release/permify/targe?include_prereleases&style=for-the-badge" alt="Targe Release" /></a>&nbsp;
    <a href="https://img.shields.io/github/commit-activity/m/Permify/targe?style=for-the-badge" target="_blank"><img src="https://img.shields.io/github/commit-activity/m/Permify/targe?style=for-the-badge" alt="Targe Commit Activity" /></a>&nbsp;
    <a href="https://img.shields.io/github/actions/workflow/status/Permify/targe/release.yml?style=for-the-badge" target="_blank"><img src="https://img.shields.io/github/actions/workflow/status/Permify/targe/release.yml?style=for-the-badge" alt="GitHub Workflow Status" /></a>&nbsp;
</p>  

![Targe Demo](assets/images/targe.gif)

## How it Works?

1. Configure your cloud credentials to enable Targe to access resources in your infrastructure. Currently, Targe supports only AWS.
2. Start an access flow or use AI to create an access command to fulfill an access request.
3. Preview the access action and complete the access request.

### Create an Access Command with AI 

Describe the access action you want to perform. For example, "give S3 read-only access to user Omer." 

Targe analyzes the request and generates the necessary access command using AI.

![targe-ai-flow](https://github.com/user-attachments/assets/0a2ea874-b6b6-47ec-b792-1602137f23e7)

### Start an Access Flow Manually

You can also manually start any flow to complete an access action. 

There are three main flows:
   - `~ % targe aws users`  | Grant or revoke access to/from a user.
   - `~ % targe aws groups` | Attach or detach a policy to/from a group.
   - `~ % targe aws roles`  | Attach or detach a policy to/from a role.

Let's repeat the example above of granting s3 read-only access to user Omer.

We will use following command to start **user** flow: `~ % targe aws users`.

The user access flow begins by listing the users in the system. Select the user to take action on.

![select-user](https://github.com/user-attachments/assets/d99327e8-3c74-42b4-9615-2afe6f0bde0b)

After selecting the user, choose the operation to perform. Let’s attach a policy to user Omer.

![select-operation](https://github.com/user-attachments/assets/bfa67375-0cd1-4bcf-9d73-dcb1a7a88dc4)

In the next step, select the policy you want to attach. You can use "filters" in each section to search what you need.

![select-policy](https://github.com/user-attachments/assets/af918b77-7e45-4c43-9d4b-f8971b1ece47)

Finally, preview the access action.

![preview-access-action](https://github.com/user-attachments/assets/d843bd92-db6d-4907-ab39-0344e4986da8)

## Installation Steps

1. **Install Targe CLI:**
   ```shell
   brew tap permify/tap-targe
   brew install targe
   ```

2. **Set Up AWS Credentials:**

   Targe requires AWS credentials to be configured in the file `~/.aws/credentials`. Follow these steps:

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

4. **Configure OpenAI API Key**
   
   Run the following command to configure your OpenAI API Key:
   ```shell
   targe config set openai_api_key [your_api_key]
   ```

5. **Set the Default Region (Optional):**

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

If you like Targe, please consider giving us a :star:

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
