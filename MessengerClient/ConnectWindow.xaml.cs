using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Navigation;
using System.Windows.Shapes;
using RestSharp;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;

namespace MessengerClient
{
    /// <summary>
    /// Interaction logic for MainWindow.xaml
    /// </summary>
    public partial class ConnectWindow : Window
    {
        public string URL;
        public string Header = "Bearer ";
        public string Name;

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


        public ConnectWindow()
        {
            InitializeComponent();
            LoginButton.Visibility = Visibility.Hidden;
            RegisterButton.Visibility = Visibility.Hidden;
        }

        private void TextBox_TextChanged(object sender, TextChangedEventArgs e)
        {
            
        }

        private void Button_Click(object sender, RoutedEventArgs e)
        {
            string IP = Text1.Text;
            string Port = Text2.Text;
            if (Port == "" || Port == "Порт (8080)")
            {
                Port = "8080";
            }
            URL = "http://"+ IP + ":" + Port;
            try
            {
                var client = new RestClient(URL);
                var request = new RestRequest("", Method.GET);
                var getresult = client.Execute(request);
                if (getresult.StatusCode.ToString() == "OK")
                {
                    dynamic json = JObject.Parse(getresult.Content);
                    Name = json.message;
                    Text1.Text = "E-mail";
                    Text1.Tag = "E-mail";
                    Text2.Text = "Пароль";
                    Text2.Tag = "Пароль";
                    Connect.Visibility = Visibility.Hidden;
                    LoginButton.Visibility = Visibility.Visible;
                    RegisterButton.Visibility = Visibility.Visible;
                    Welcome_Label.Content = "Добро пожаловать на сервер:";
                    Server_name.Content = Name;
                }
                else
                {
                    MessageBox.Show("Не удалось подключиться к серверу", "Ошибка", MessageBoxButton.OK, MessageBoxImage.Stop);
                }
            }
            catch {
                MessageBox.Show("Введены неверные данные", "Ошибка", MessageBoxButton.OK, MessageBoxImage.Stop);
            }
        }

        private void ListBox_SelectionChanged(object sender, SelectionChangedEventArgs e)
        {

        }

        private void LoginButton_Click(object sender, RoutedEventArgs e)
        {
            var client = new RestClient(URL);
            var request = new RestRequest("/api/user/login", Method.POST, RestSharp.DataFormat.Json);
            request.AddJsonBody(new { email = Text1.Text, password = Text2.Text });
            var getresult = client.Execute(request);
            if (getresult.StatusCode.ToString() == "OK") {
                dynamic json = JObject.Parse(getresult.Content);
                Header += json.account.token;
                Chat chat = new Chat(URL, Header,Name);
                chat.Show();
                this.Close();
            }
            else
            {
                dynamic json = JObject.Parse(getresult.Content);
                if(json.message == "USER_NOT_FOUND") MessageBox.Show("Пользователь не зарегистрирован", "Ошибка", MessageBoxButton.OK, MessageBoxImage.Stop);
                if (json.message == "INCORRECT_PASS") MessageBox.Show("Неверный пароль", "Ошибка", MessageBoxButton.OK, MessageBoxImage.Stop);
            }
        }

        private void RegisterButton_Click(object sender, RoutedEventArgs e)
        {
            var client = new RestClient(URL);
            var request = new RestRequest("/api/user/create", Method.POST, RestSharp.DataFormat.Json);
            request.AddJsonBody(new { email = Text1.Text, password = Text2.Text });
            var getresult = client.Execute(request);
            if (getresult.StatusCode.ToString() == "OK")
            {
                MessageBox.Show("Теперь вы можете войти в свой аккаунт", "Успех", MessageBoxButton.OK, MessageBoxImage.Information);
            }
            else
            {
                dynamic json = JObject.Parse(getresult.Content);
                if (json.message == "ALREADY_USED") MessageBox.Show("Пользователь уже зарегистрирован", "Ошибка", MessageBoxButton.OK, MessageBoxImage.Stop);
                if (json.message == "LEN") MessageBox.Show("Пароль должен быть длиной не меньше 6 символов", "Ошибка", MessageBoxButton.OK, MessageBoxImage.Stop);
                if (json.message == "NOEMAIL") MessageBox.Show("Неверный адрес электронной почты", "Ошибка", MessageBoxButton.OK, MessageBoxImage.Stop);
            }
        }
    }
}
