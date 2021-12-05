using Newtonsoft.Json.Linq;
using RestSharp;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using System.Timers;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Shapes;
namespace System.Windows.Controls
{
    public static class MyExt
    {
        public static void PerformClick(this Button btn)
        {
            btn.RaiseEvent(new RoutedEventArgs(Button.ClickEvent));
        }
    }
    public static class Ext2
    {
        public static void AppendText(this RichTextBox box, string text, SolidColorBrush brush)
        {
            TextRange tr = new TextRange(box.Document.ContentEnd, box.Document.ContentEnd);
            tr.Text = text;
            tr.ApplyPropertyValue(TextElement.ForegroundProperty, brush);
        }
    }
}

namespace MessengerClient
{
    /// <summary>
    /// Логика взаимодействия для Chat.xaml
    /// </summary>
    ///

    public partial class Chat : Window
    {
        string URL;
        string Header;
        string Name;
        string CurUser;
        int err;
        System.Windows.Threading.DispatcherTimer timer = new System.Windows.Threading.DispatcherTimer();
        public void RemoveText(object sender, EventArgs e)
        {
            TextBox instance = (TextBox)sender;
            if (instance.Text == instance.Tag.ToString())
                instance.Text = "";
        }

        public void AddText(object sender, EventArgs e)
        {
            TextBox instance = (TextBox)sender;
            if (string.IsNullOrWhiteSpace(instance.Text))
                instance.Text = instance.Tag.ToString();
        }
        public Chat(string A, string B, string C) {
            URL = A;
            Header = B;
            Name = C;
            InitializeComponent();
            this.Title = "Kursach Messenger - " + Name + "(" + URL + ")";
            var usrclient = new RestClient(URL);
            var usrrequest = new RestRequest("/api/user", Method.GET);
            usrrequest.AddHeader("Authorization", Header);
            var getuser = usrclient.Execute(usrrequest);
            dynamic jsonusr = JObject.Parse(getuser.Content);
            var CurUseremail = jsonusr.account.email;
            CurUser = CurUseremail.ToString();
            var client = new RestClient(URL);
            var request = new RestRequest("/api/messages/get", Method.GET);
            request.AddHeader("Authorization",Header);
            var getresult = client.Execute(request);
            if (getresult.StatusCode.ToString() == "OK")
            {
                dynamic json = JObject.Parse(getresult.Content);
                string amount = json.message;
                int end = int.Parse(amount);
                for(int i = 0; i < end; i++)
                {
                    string indx = i.ToString();
                    dynamic jsoni = json[i.ToString()];
                    DateTime datetime = jsoni["CreatedAt"];
                    string time = datetime.ToString("MM/dd/yyyy HH:mm:ss");
                    ChatWindow.AppendText(time + " ", Brushes.Gray);
                    if (jsoni["name"] == CurUser)
                    {
                        string k = jsoni["name"].ToString();
                        ChatWindow.AppendText(k, Brushes.BlueViolet);
                    }
                    else
                    {
                        string k = jsoni["name"].ToString();
                        ChatWindow.AppendText(k, Brushes.IndianRed);
                    }
                    
                    ChatWindow.AppendText(": ",Brushes.Gray);
                    string msg = jsoni["msg"].ToString();
                    ChatWindow.AppendText(msg, Brushes.Black);
                    ChatWindow.AppendText("\n");
                }
                ChatWindow.ScrollToEnd();
            }
            timer.Tick += new EventHandler(timerTick);
            timer.Interval = new TimeSpan(0, 0, 0,0,500);
            timer.Start();
            err = 0;
        }
        private void timerTick(object sender, EventArgs e)
        {
            var client = new RestClient(URL);
            var request = new RestRequest("/api/messages/get/unread", Method.GET);
            request.AddHeader("Authorization", Header);
            try
            {
                var getresult = client.Execute(request);
                dynamic json = JObject.Parse(getresult.Content);
                if (getresult.StatusCode.ToString() == "OK" && json.message != "NO_NEW_MSG")
                {
                    string amount = json.message;
                    int end = int.Parse(amount);
                    for (int i = 0; i < end; i++)
                    {

                        string indx = i.ToString();
                        dynamic jsoni = json[i.ToString()];
                        DateTime datetime = jsoni["CreatedAt"];
                        string time = datetime.ToString("MM/dd/yyyy HH:mm:ss");
                        ChatWindow.AppendText(time + " ", Brushes.Gray);
                        if (jsoni["name"] == CurUser)
                        {
                            string k = jsoni["name"].ToString();
                            ChatWindow.AppendText(k, Brushes.BlueViolet);
                        }
                        else
                        {
                            string k = jsoni["name"].ToString();
                            ChatWindow.AppendText(k, Brushes.IndianRed);
                        }

                        string msg = jsoni["msg"].ToString();
                        ChatWindow.AppendText(": ", Brushes.Gray);
                        ChatWindow.AppendText(msg, Brushes.Black);
                        ChatWindow.AppendText("\n");
                    }
                }
            }
            catch
            {
                MessageBox.Show("Соединение с сервером прервано!", "Ошибка", MessageBoxButton.OK, MessageBoxImage.Stop);
                ConnectWindow connect = new ConnectWindow();
                connect.Show();
                this.Close();
                timer.Stop();
            }

        }
        private void Button_Click(object sender, RoutedEventArgs e)
        {
            var client = new RestClient(URL);
            var request = new RestRequest("/api/messages/send", Method.POST);
            if(SendMsgBox.Text == SendMsgBox.Tag)
            {
                request.AddJsonBody(new { msg = "" });
                return;
            }
            else request.AddJsonBody(new { msg = SendMsgBox.Text });
            request.AddHeader("Authorization", Header);
            var getresult = client.Execute(request);
            if (getresult.StatusCode.ToString() != "OK")
            {
                dynamic json = JObject.Parse(getresult.Content);
                if (json.message == "EMPTY_MSG") MessageBox.Show("Сообщение не должно быть пустым!", "Ошибка", MessageBoxButton.OK, MessageBoxImage.Error);
            }
            SendMsgBox.Text = "";
            SendMsgBox.Tag = "Введите сюда сообщение...";
        }
        private void RichTextBox_TextChanged(object sender, TextChangedEventArgs e)
        {
            ChatWindow.ScrollToEnd();
        }

        private void SendMsgBox_TextChanged(object sender, TextChangedEventArgs e)
        {

        }
        private void UserControl_Loaded(object sender, RoutedEventArgs e)
        {
            var window = Window.GetWindow(this);
            window.KeyDown += HandleKeyPress;
        }


        private void HandleKeyPress(object sender, KeyEventArgs e)
        {
            if (e.Key == Key.Enter)
            {
                Send.PerformClick();
            }
        }
    }
}
