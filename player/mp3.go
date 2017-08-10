// Copyright 2017 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package player

import (
	"io"
	"os"

	"github.com/hajimehoshi/oto"

	"github.com/hajimehoshi/go-mp3"
)

func Play(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	d,err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}
	defer d.Close()
	/*var frame mp3.Frame
	skip := 0
	d.Decode(&frame, &skip)*/

	p, err := oto.NewPlayer(44100, 2, 2, 8192*4)
	if err != nil {
		return err
	}
	defer p.Close()


	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil
}
