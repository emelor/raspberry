en paj har kanske en funktion 
Id() string

och en hub har en 
map[string]Pi

så att hubbens 
RegisterPi(p Pi) 

blir något i stil med
self.pis[p.Id()] = pi

så i test.go
p := pi.NewPi( asdafds )
h := hub.NewHub( asdfasdfafds )
h.start()
p.ConnectTo(h)


func (self *Pi) ConnectTo(h common.Hub) {
  h.Register(self)
  self.startTimers()
}
func (self *Hub) Register(p common.Pi) {
  self.pis[p.Id()] = p
}
kanske skall test.go ha en "h.Start()" mellan p := ... och p.ConnectTo...





if blabla {
	if gris {
		kasta		
	}
} else do bläblä