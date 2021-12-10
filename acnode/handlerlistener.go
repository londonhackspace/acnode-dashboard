package acnode

import "github.com/rs/zerolog/log"

type HandlerListener struct {
	handler *ACNodeHandler
	name    string

	nodeAdded   chan ACNode
	nodeChanged chan ACNode

	handlerChangeAdded   chan func(node ACNode)
	handlerChangeChanged chan func(node ACNode)

	onNodeAdded   func(node ACNode)
	onNodeChanged func(node ACNode)
}

func CreateHandlerChangeListener(handler *ACNodeHandler, name string) *HandlerListener {
	listener := HandlerListener{
		handler:              handler,
		name:                 name,
		nodeAdded:            make(chan ACNode),
		nodeChanged:          make(chan ACNode),
		handlerChangeAdded:   make(chan func(ACNode)),
		handlerChangeChanged: make(chan func(ACNode)),
		onNodeAdded:          nil,
		onNodeChanged:        nil,
	}

	handler.AddListener(&listener)
	go listener.runACNodeHandlerListener()

	return &listener
}

func (sub *HandlerListener) runACNodeHandlerListener() {
	sub.handler.AddListener(sub)
	defer func() {
		sub.handler.RemoveListener(sub)
		// other channels will be closed by the listener removal
		close(sub.handlerChangeAdded)
		close(sub.handlerChangeChanged)
	}()
	log.Info().Str("Listener", sub.name).Msg("ACNode Listener running")
	for {
		select {
		case node, ok := <-sub.nodeAdded:
			if !ok {
				break
			}
			if sub.onNodeAdded != nil {
				sub.onNodeAdded(node)
			}
		case node, ok := <-sub.nodeChanged:
			if !ok {
				break
			}
			if sub.onNodeChanged != nil {
				sub.onNodeChanged(node)
			}
		case f := <-sub.handlerChangeAdded:
			sub.onNodeAdded = f
		case f := <-sub.handlerChangeChanged:
			sub.onNodeChanged = f
		}
	}
	log.Info().Str("Listener", sub.name).Msg("ACNode Listener disconnected")
}

func (sub *HandlerListener) SetOnNodeChangedHandler(h func(ACNode)) {
	sub.handlerChangeChanged <- h
}

func (sub *HandlerListener) SetOnNodeAddedHandler(h func(ACNode)) {
	sub.handlerChangeAdded <- h
}

func (sub *HandlerListener) Disconnect() {
	sub.handler.RemoveListener(sub)
}
