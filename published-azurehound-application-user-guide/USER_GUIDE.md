**Steps to Deploy the SpecterOps AzureHound Managed Application from the Azure Marketplace**
1. Log in to the Azure Portal
2. In the Azure Portal, search and select Marketplace.
3. Use the search bar to find **SpecterOps AzureHound**.
<img width="773" height="660" alt="image10" src="https://github.com/user-attachments/assets/4e040def-01c8-4643-88f0-a7821ba55808" />

4. From the search results, click on SpecterOps AzureHound to open the product page.
5. Click the "Create" button to begin the deployment process.
<img width="1023" height="791" alt="image11" src="https://github.com/user-attachments/assets/10fb6bfb-8c31-4682-9f42-2eeee3f873ce" />



6. Configure Basic Settings :
- Choose the **Subscription** under which to deploy the application.
- Select or create a **Resource Group**.
- Enter a **Region** for deployment.
- Provide a name for your **Managed Application**, then click Next.
<img width="832" height="869" alt="image5" src="https://github.com/user-attachments/assets/786de45b-7506-4a2e-b0bb-dae1430404fe" />


**Steps to Register an Application in Microsoft Entra ID**
- Sign in to the **Microsoft Entra** Admin Center.
- Navigate to **Entra ID** > **App registrations**.
- Click - **New registration**.
- Provide a meaningful Name for the Application (e.g., azurehound-client-app).
- Under Supported account types, select - **Accounts in this organizational directory only**.
- Click Register to create the Application.
- Once registration is complete, you'll be redirected to the application’s Overview page.
- Copy and save the **Application (client) ID** — you'll need it during deployment.

 **Create a new Application Secret**
- Select App registrations and select your application from the list.
- Select **Certificates & secrets**.
- Select Client secrets, and then Select New client secret.
- Provide a description of the secret, and a duration.
- Select **Add**.

**To obtain the BloodHound Token ID and Token Secret, follow these steps:**
- Log in to your **BloodHound instance**.
- Navigate to **Administration** > **Manage Clients**.
- Click Create Client to generate a new Managed Client.
- Once created, copy the **Token ID** and **Token Secret** — these values will be used in the application deployment parameters.

<img width="801" height="458" alt="image1" src="https://github.com/user-attachments/assets/7507863d-23f6-4ea6-861d-78c493a83837" />


7. Fill in the required AzureHound Config Params:
- **Azure Tenant ID** - Your Azure Tenant ID.
- **Azure Application ID** - Register an application in Microsoft Entra ID, and grant it the Directory.Read.All, RoleManagement.Read.All API permissions and admin consent.

<img width="1071" height="286" alt="image7" src="https://github.com/user-attachments/assets/aaca6471-2935-4725-afaf-fb057cd6fb6a" />


- **Azure Secret ID** - Create a Client Secret for the registered app, and enter the secret value (not the ID).
- **BloodHound Instance Domain** - Enter your BloodHound instance domain name
- **BloodHound Token ID** - Enter the Managed Client Token ID.
- **BloodHound Token Secret** -Enter the Managed Client Token Secret.
- **Azure Function Package** - Enter URI to access Azure Function Package
  https://saazurehounddev.blob.core.windows.net/azurefunction/containerRestartFunction.zip 

8. Click Next, then Review + Create.
After validation, click Create to begin deployment.


**Start a job in Bloodhound**

After creating a Client in the Manage Clients section of BloodHound:
- Locate your client in the list.
- Click the menu icon (three horizontal lines) on the right side of the client row.
- Select On Demand Scan and click Run to start the job immediately.
- Optional - Schedule a Job (If you'd like AzureHound to run on a regular schedule)
Click Edit Client > Configure the Collection Schedule based on your preferred timing and frequency > Save the changes to apply the schedule.

<img width="348" height="891" alt="image3" src="https://github.com/user-attachments/assets/61848c27-9d3d-4a88-93e0-4dd4e389da36" />
<img width="348" height="274" alt="image2" src="https://github.com/user-attachments/assets/61526c96-fd95-4f0f-84d0-d84dd25161e7" />

After the deployment is finished, you can check your managed application's status.
Navigate to the resource group you selected during the deployment. Under the Overview tab, you will find your deployed Managed Application listed among the resources.

<img width="1334" height="533" alt="image6" src="https://github.com/user-attachments/assets/b0d5fb11-bcaa-47a8-a993-eceb63c7b3b7" />

Click on your **deployed Managed Application**
<img width="1187" height="658" alt="new-image-" src="https://github.com/user-attachments/assets/79b35322-4e80-434c-a1a9-1a8541619dc1" />
Click on the Managed resource group, and you can see the resources deployed.

**View Logs from the Deployed Container App**
To monitor and troubleshoot your AzureHound deployment, you can access real-time logs from the container app:
- In the Managed resource group, in Resources, search for the container app and open the Container App resource.
  
<img width="1080" height="323" alt="image4" src="https://github.com/user-attachments/assets/5f5675e1-7bb1-40a9-8c45-b0b51f3cf1d3" />

- In the left-hand search bar within the Container App blade, type Log Stream and select it from the options.
- Set the Display to Real-Time.
- Under Category, select Application to view logs generated by the AzureHound application.
<img width="977" height="465" alt="image8" src="https://github.com/user-attachments/assets/54d22ee8-d35d-4463-88d6-8bed15c4255b" />

By following the steps outlined in this guide, you can successfully deploy and configure AzureHound as a Managed Application in Microsoft Azure. This streamlined approach ensures minimal manual setup, secure integration with Azure services, and seamless visualization of Azure data within the BloodHound platform.
