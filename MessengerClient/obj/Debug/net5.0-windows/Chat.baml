   M S B A M L   `   `   `  ����J  FMessengerClient, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null1clr-namespace:MessengerClientMessengerClient  5   	       MessengerClient.Chat   5   	   W SPresentationCore, Version=5.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35\ XPresentationFramework, Version=5.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35e aSystem.Windows.Controls.Ribbon, Version=5.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089R NWindowsBase, Version=5.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35F 9http://schemas.microsoft.com/winfx/2006/xaml/presentation     5   	   R NSystem.Xaml, Version=5.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089:x,http://schemas.microsoft.com/winfx/2006/xaml     5   	   @d2http://schemas.microsoft.com/expression/blend/2008     5   	   Jmc;http://schemas.openxmlformats.org/markup-compatibility/2006     5   	   /localclr-namespace:MessengerClient     5   	     � Title$  Kursach Messenger��5   	   $	��479��6#   $	��828��60   .��5	      � $��0,10,0,4q�6    � RowDefinitions 5
   
   �� 5      $
��227* �6   �� 5      $
��232* �6   .+�5   
   >� -   5        5      �� 5   
   -   5         ��Name$ 
ChatWindow��5      $��10,10,10,50q�6*    �� 
IsReadOnly$
 True��6?   $��2��6w   . �5      /� .T�5      J� .*�6   .*�5      �� $� ��6   5      5      5      �� 5   
   -   5         ��Name$ 
SendMsgBox��5      $��10,0,127,10q�6&   $3��-Введите сюда сообщение...��6;    �� TextWrapping$
 Wrap=�6\   $��32��6p   $��Bottom=�6|   
 � Tag$3 -Введите сюда сообщение...��6�   $��1��6  �� 5   
   -   5         ��Name$
 Send��5      $��Отправить��6   $��	0,0,10,10q�63   $��32��6[   $��Bottom=�6g   $��1��6�   $��Right=�6�   $��94��6�   5      5      