package utils

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)


//           _     _   _        
//          | |   (_) | |       
//   _   _  | |_   _  | |  ___  
//  | | | | | __| | | | | / __| 
//  | |_| | \ |_  | | | | \__ \ 
//   \__,_|  \__| |_| |_| |___/ 

// ConstructorSrv заполняет нашу структуру
func (s *Server) ConstructorSrv(port string, maxConnections int) error {
	srv, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	// устанавливаем настройки нашего сервера
	s.Server = srv
	s.MaxConnections = maxConnections
	s.Connections = make(map[net.Conn]string, maxConnections)
	s.UsedNames = make(map[string]bool, maxConnections)
	return nil
}

//CanConnect проверяет может ли юзер войти в систему
func (s *Server) CanConnect(conn net.Conn) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// проверяем может ли юзер добавится к чату
	return !(s.MaxConnections != 0 && len(s.Connections) >= s.MaxConnections)
}

// ConnectMessenger добавляет на сервер юзера с именем
func (s *Server) ConnectMessenger(conn net.Conn) {
	// проверяем может ли юзер добавится к чату
	if !s.CanConnect(conn) {
		// если нет то уведомляем его
		fmt.Fprint(conn, "The room is full, please try again later.")
		// закрываем соединение
		conn.Close()
		return
	}
	// приветсвуем нашего пользователя
	// Fprint принимает в первый аргумент куда будет писать сообщение
	fmt.Fprint(conn, WelcomMessage)
	// создём новый Reader
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Printf("ConnectMessenger: %v\n", err.Error())
		return
	}
	// получаем имя
	name = strings.Replace(name, "\n", "", 1)
	// добавляем юзера в сервер
	err = s.addConnection(conn, name)
	if err != nil {
		fmt.Fprint(conn, err.Error())
		conn.Close()
		return
	}
	// начинаем общение
	s.StartChating(conn)
	// удаляем соединение если пользователь нажал на ctrl c
	s.RemoveConnection(conn)
}
//StartChating функция которая выгружает все смс отправленные до пользователя и добавляет его самого в чат
func (s *Server) StartChating(conn net.Conn) {
	// загружаем все письма предыдущих пользователей к новому пользователю
	s.DownloadMessages(conn)
	// готовим сообщение о том что новый пользователь присоединился
	sms := getFormattedMessage(s, conn, "", ModeJoinChat)
	// отправляем в общий чат
	s.SendMessage(conn, sms)
	// запускаем бесконечный цикл который слушает что напишет пользователь постоянно
	for {
		// создаёт поле для чтения
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}
		// готовит сообщение к отправке
		message = getFormattedMessage(s, conn, message, ModeSendMessage)
		// отправляем это сообщение
		s.SendMessage(conn, message)
		// сохраняет это сообщение
		s.SavedMessage(message)
	}
	// если он нажал ctrl c то срабатывает if внутри цикла и мы переходим к это строчке кода которое готовит сообщение о том что пользователь вышел
	message := getFormattedMessage(s, conn ,"", ModeLeftChat)
	// отправляем это сообщение в группу 
	s.SendMessage(conn, message)
}

//SavedMessage сохраняет письмо пользователя в общую структуру сервера
func (s *Server) SavedMessage(message string) {
	s.mutex.Lock()
	s.AllMessages = append(s.AllMessages, message)
	s.mutex.Unlock()
}
//SendMessage отправляет всем юзерам на сервере
func (s *Server) SendMessage(conn net.Conn, message string) {
	// проверяем на пустоту
	if message == "" {
		fmt.Fprintf(conn, PatternSending, time.Now().Format(TimeDefault), s.Connections[conn])
		return
	}
	// отправка сообщений
	time := time.Now().Format(TimeDefault)
	sendMessage := fmt.Sprintf("%s!%s\n%s", ColorYellow, ColorReset, message)
	s.mutex.Lock()
	for con := range s.Connections {
		if con != conn {
			fmt.Fprint(con, sendMessage)
		}
		fmt.Fprintf(con, PatternSending, time, s.Connections[con])
	}
	s.mutex.Unlock()
}

func getFormattedMessage(srv *Server, conn net.Conn, message string, mode int) string {
	srv.mutex.Lock()
	// лочим поток

	// получаем имя пользователя
	name := srv.Connections[conn]
	srv.mutex.Unlock()
	// Change Message
	switch mode {
	case ModeSendMessage:
		if message == "\n" {
			return ""
		}
		time := time.Now().Format(TimeDefault)
		// в зависимости от mode отправляем сообщение это вывод обычного сообщения
		message = fmt.Sprintf(PatternMessage, time, name, message)
	case ModeJoinChat:
		// уведомляем о соединении нового юзера
		message = fmt.Sprintf(ColorYellow+PatternJoinChat+ColorReset, name)
	case ModeLeftChat:
		// уведомляем о выходе пользователя
		message = fmt.Sprintf(ColorYellow+PatternLeftChat+ColorReset, name)
	}
	// отправляем готовое сообщение
	return message
}

func (s *Server) addConnection(conn net.Conn, name string) error {
	// проверяем имя на пустоту
	if name == "" {
		return errors.New("Name cant be empty")
	}
	// блокируем доступ другим горутинам
	s.mutex.Lock()
	//  после завершения функции открываем обратно
	defer s.mutex.Unlock()
	// проверяем может ли юзер зайти в систему
	if s.MaxConnections != 0 && len(s.Connections) >= s.MaxConnections {
		return fmt.Errorf("the room is full (%v)", conn.RemoteAddr())
	} else if s.UsedNames[name] {
		//  проверка есть ли такое имя уже в системе
		return fmt.Errorf("name '%s' is exist [%v] ", name, conn.RemoteAddr())
	}
	// добавляем имя в общую структуру
	s.UsedNames[name] = true
	// добавляем в мапу коннект этого юзера и его имя в значении
	s.Connections[conn] = name
	// уведомляем пользователей о вхождении в систему
	log.Printf("Connect %v ", conn.RemoteAddr())
	return nil
}

// DownloadMessages загружает все письма пользователей в интерфейс нового пользователя
func (s *Server) DownloadMessages(conn net.Conn) {
	for _, sms := range s.AllMessages {
		fmt.Fprintf(conn, sms)
	}
}


// RemoveConnection - Removing Connection from s.Connections (safe)
func (s *Server) RemoveConnection(conn net.Conn) {
	// лочим поток от других горутин
	s.mutex.Lock()
	// удаляем указанный элемент с мапы
	delete(s.UsedNames, s.Connections[conn])
	// удаляем указанный элемент с другой мапы
	delete(s.Connections, conn)
	// уведомляем пользователей о том что юзер вышел из чата
	log.Printf("Connect %v was left", conn.RemoteAddr())
	// открываем соединение для других горутин
	s.mutex.Unlock()
}
