package shell

type directory struct {
	name string
	parent *directory
	child []*directory
	action action
}

func (d *directory) addChild(child *directory){
	d.child = append(d.child, child)
	child.parent = d
}

func newDirectory(name string, actions ...map[string]action) *directory{
	return &directory{
		name: name,
		action: func(s *shell, inputs []byte) (string, error) {
			return distribute(s, inputs, actions...)
		},
	}
}

var root = newDirectory("/", rootActions, baseActions)
var posts = newDirectory("posts", postsActions, baseActions)
var followers = newDirectory("followers", followersActions, baseActions)

func init(){
	root.addChild(posts)
	root.addChild(followers)
}
