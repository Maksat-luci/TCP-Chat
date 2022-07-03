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
	// 
	name = strings.Replace(name, "\n", "", 1)
	err = s.addConnection(conn, name)
	if err != nil {
		fmt.Fprint(conn, err.Error())
		conn.Close()
		return
	}
	s.StartChating(conn)
	s.RemoveConnection(conn)
}
//StartChating функция которая выгружает все смс отправленные до пользователя и добавляет его самого в чат
func (s *Server) StartChating(conn net.Conn) {
	s.DownloadMessages(conn)
	sms := getFormattedMessage(s, conn, "", ModeJoinChat)
	s.SendMessage(conn, sms)
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}
		message = getFormattedMessage(s, conn, message, ModeSendMessage)
		s.SendMessage(conn, message)
		s.SavedMessage(message)
	}
	message := getFormattedMessage(s, conn ,"", ModeLeftChat)
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
	name := srv.Connections[conn]
	srv.mutex.Unlock()
	// Change Message
	switch mode {
	case ModeSendMessage:
		if message == "\n" {
			return ""
		}
		time := time.Now().Format(TimeDefault)
		message = fmt.Sprintf(PatternMessage, time, name, message)
	case ModeJoinChat:
		message = fmt.Sprintf(ColorYellow+PatternJoinChat+ColorReset, name)
	case ModeLeftChat:
		message = fmt.Sprintf(ColorYellow+PatternLeftChat+ColorReset, name)
	}
	return message
}

func (s *Server) addConnection(conn net.Conn, name string) error {
	if name == "" {
		return errors.New("Name cant be empty")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.MaxConnections != 0 && len(s.Connections) >= s.MaxConnections {
		return fmt.Errorf("the room is full (%v)", conn.RemoteAddr())
	} else if s.UsedNames[name] {
		return fmt.Errorf("name '%s' is exist [%v] ", name, conn.RemoteAddr())
	}
	s.UsedNames[name] = true
	s.Connections[conn] = name
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
	s.mutex.Lock()
	// удаляем указанный элемент с мапы
	delete(s.UsedNames, s.Connections[conn])
	delete(s.Connections, conn)
	// уведомляем пользователей о том что юзер вышел из чата
	log.Printf("Connect %v was left", conn.RemoteAddr())
	s.mutex.Unlock()
}
