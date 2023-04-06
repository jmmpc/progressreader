package progressreader

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"
)

const loremIpsum = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus nunc leo, rutrum a nulla sit amet, laoreet facilisis velit. Nullam id luctus tellus. Proin eu nibh non velit sollicitudin consequat. Nullam ac nunc facilisis, fermentum arcu sit amet, posuere eros. Etiam a tempus dolor, sit amet posuere mi. Maecenas nibh elit, fermentum at lacus ac, pulvinar fermentum sapien. Phasellus consectetur urna justo, non dapibus nisl congue id. Ut aliquet lacus neque, vitae pharetra metus fermentum a. Sed vehicula mollis mollis. Proin varius risus quis lacus facilisis, a posuere lectus hendrerit. Quisque metus risus, sollicitudin vitae sodales a, tempor in velit. Sed varius leo tortor, ut finibus nisi venenatis et. Sed iaculis justo nisi, eget volutpat nisi mattis sit amet. Morbi vel iaculis ante. Ut ultrices, elit eget ullamcorper gravida, lectus justo ornare sapien, vel facilisis orci leo ut nibh. Praesent turpis eros, blandit eget odio convallis, dictum hendrerit nisl.

Duis fringilla pharetra pulvinar. Cras id leo eu turpis suscipit viverra. Sed feugiat cursus purus, sit amet semper orci scelerisque id. Nullam malesuada, nibh non gravida posuere, est leo venenatis ex, at malesuada ante diam vel turpis. Nunc id sodales enim, sed bibendum nisi. Proin condimentum, arcu ac viverra ullamcorper, purus dui vehicula erat, id faucibus est tellus et lorem. Fusce rutrum venenatis dapibus. Cras vestibulum faucibus commodo. Nunc dui magna, convallis sit amet orci eget, egestas rhoncus tortor. Mauris fermentum mi ut neque luctus egestas sed sagittis dolor. Donec venenatis, eros ut pharetra rhoncus, eros lacus rutrum ex, vel malesuada sapien lectus eu tellus. Sed sit amet rhoncus ante. Ut ultrices aliquam suscipit. Suspendisse laoreet pretium tincidunt. Suspendisse eu consequat lectus.

Sed at nisl tortor. Sed bibendum turpis at metus semper, id imperdiet sem euismod. Praesent augue mauris, bibendum eget vulputate sed, ornare id ex. Nullam euismod nisl at nisl tristique maximus. Quisque tincidunt fringilla rutrum. Vestibulum sed eros nec leo elementum posuere nec a mauris. Sed nulla elit, lobortis nec blandit non, vehicula sit amet enim. Aliquam erat volutpat.

Aenean et ullamcorper elit. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris interdum ex et mi tempus dictum. Phasellus eget mauris tortor. Donec tincidunt ac nisi at venenatis. Pellentesque commodo tincidunt libero non dignissim. Etiam interdum rhoncus consectetur. Nam eget egestas ligula. Suspendisse ac eros sit amet justo cursus vestibulum. Curabitur in ex condimentum, luctus massa vitae, sodales sapien. Vivamus ullamcorper nunc sit amet sollicitudin tempus.

Mauris nec euismod urna. Cras faucibus dignissim euismod. Morbi rutrum lacinia hendrerit. Proin fermentum erat et nibh vestibulum egestas. Aenean sit amet accumsan urna, id auctor ipsum. Sed eget nulla rhoncus, imperdiet enim ut, lobortis massa. Vestibulum ut blandit velit, eget finibus enim.`

func TestProgressReader(t *testing.T) {
	var out bytes.Buffer

	in := New(bytes.NewBufferString(loremIpsum))

	n, err := io.Copy(&out, in)
	if err != nil {
		t.Fatal(err)
	}

	if n != int64(len(loremIpsum)) {
		t.Fatalf("wrong number of bytes read: expected %d, has %d\n", len(loremIpsum), n)
	}

	if out.String() != loremIpsum {
		t.Fatal("input is not equal to output")
	}
}

type FakeReader struct {
}

func (r *FakeReader) Read(b []byte) (int, error) {
	time.Sleep(20 * time.Millisecond)
	return 512, nil
}

func TestProgressReaderWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	in := WithContext(ctx, new(FakeReader))
	out := new(bytes.Buffer)

	done := make(chan error)

	go func() {
		n, err := io.Copy(out, in)
		if n == 0 {
			t.Fatal("Read bytes number must be greater than 0")
		}
		done <- err
	}()

	for {
		select {
		case <-time.Tick(15 * time.Millisecond):
			_ = in.Total()
		case err := <-done:
			if err != context.DeadlineExceeded {
				t.Fatalf("wrong error value: expected %v, has %v\n", context.DeadlineExceeded, err)
			}
			return
		}
	}
}
