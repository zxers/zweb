package zweb

type Middleware func(next HandleFunc) HandleFunc