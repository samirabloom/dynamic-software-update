import io.netty.bootstrap.ServerBootstrap;
import io.netty.buffer.Unpooled;
import io.netty.channel.*;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.SocketChannel;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import io.netty.handler.codec.http.*;
import io.netty.handler.logging.LoggingHandler;

import java.nio.charset.Charset;
import java.util.concurrent.TimeUnit;

import static io.netty.handler.codec.http.HttpHeaders.Names.*;
import static io.netty.handler.codec.http.HttpResponseStatus.OK;
import static io.netty.handler.codec.http.HttpVersion.HTTP_1_1;

public class TestServer {

    public static void main(final String[] args) throws InterruptedException {
        System.out.println("STARTING SERVER FOR HTTP ON PORT: " + args[0]);
        NioEventLoopGroup bossGroup = new NioEventLoopGroup();
        NioEventLoopGroup workerGroup = new NioEventLoopGroup();
        try {
            new ServerBootstrap()
                    .group(bossGroup, workerGroup)
                    .channel(NioServerSocketChannel.class)
                    .childHandler(new ChannelInitializer<SocketChannel>() {
                        @Override
                        public void initChannel(SocketChannel ch) throws Exception {
                            ChannelPipeline pipeline = ch.pipeline();

                            pipeline.addLast("logger", new LoggingHandler("TEST_SERVER"));
                            pipeline.addLast("http_codec", new HttpServerCodec());
                            pipeline.addLast("simple_test_handler", new TestServerHandler(Integer.parseInt(args[0])));
                        }
                    })
                    .bind(Integer.parseInt(args[0])).channel().closeFuture().sync();
        } catch (Exception e) {
            throw new RuntimeException("Exception running test server", e);
        } finally {
            bossGroup.shutdownGracefully();
            workerGroup.shutdownGracefully();
        }
    }

    private static class TestServerHandler extends SimpleChannelInboundHandler<Object> {
        int port;

        public TestServerHandler(int _port) {
            port = _port;
        }

        @Override
        public void channelReadComplete(ChannelHandlerContext ctx) {
            ctx.flush();
        }

        @Override
        public void channelRead0(ChannelHandlerContext ctx, Object msg) throws Exception {
            if (msg instanceof HttpRequest) {
                HttpRequest req = (HttpRequest) msg;
                if (req.getUri().equals("/chunked")) {
                    HttpResponse response = new DefaultHttpResponse(HTTP_1_1, OK);
                    response.headers().set(CACHE_CONTROL, "no-cache, no-store, must-revalidate");
                    response.headers().set(PRAGMA, "no-cache");
                    response.headers().set(EXPIRES, "0");
                    response.headers().set(HttpHeaders.Names.TRANSFER_ENCODING, HttpHeaders.Values.CHUNKED);
                    response.headers().set(CONNECTION, HttpHeaders.Values.KEEP_ALIVE);
                    ctx.write(response);

                    ctx.write(new DefaultHttpContent(Unpooled.wrappedBuffer(String.format("NS Port: %d\n", port).getBytes(Charset.forName("UTF-8")))));
                    TimeUnit.MILLISECONDS.sleep(50);
                    ctx.write(new DefaultHttpContent(Unpooled.wrappedBuffer("50 ms, ".getBytes(Charset.forName("UTF-8")))));
                    ctx.flush();
                    TimeUnit.MILLISECONDS.sleep(50);
                    ctx.write(new DefaultHttpContent(Unpooled.wrappedBuffer("100 ms, ".getBytes(Charset.forName("UTF-8")))));
                    ctx.flush();
                    TimeUnit.MILLISECONDS.sleep(50);
                    ctx.write(new DefaultHttpContent(Unpooled.wrappedBuffer("150 ms, ".getBytes(Charset.forName("UTF-8")))));
                    ctx.flush();
                    TimeUnit.MILLISECONDS.sleep(50);
                    ctx.write(new DefaultHttpContent(Unpooled.wrappedBuffer("200 ms, ".getBytes(Charset.forName("UTF-8")))));
                    ctx.flush();
                    TimeUnit.MILLISECONDS.sleep(50);
                    ctx.write(new DefaultHttpContent(Unpooled.wrappedBuffer("250 ms\n".getBytes(Charset.forName("UTF-8")))));
                    ctx.flush();

                    // close
                    ctx.write(new DefaultHttpContent(Unpooled.wrappedBuffer("\r\n".getBytes(Charset.forName("UTF-8")))));
                    ctx.write(new DefaultHttpContent(Unpooled.wrappedBuffer("".getBytes(Charset.forName("UTF-8")))));
                    ctx.write(new DefaultHttpContent(Unpooled.wrappedBuffer(new byte[]{})));
                    ctx.write(new DefaultHttpContent(Unpooled.EMPTY_BUFFER));
                    ctx.flush();
                } else if (req.getUri().equals("/close")) {
                    FullHttpResponse response;
                    String serverResponse = String.format("NS Port: %d\n", port);
                    TimeUnit.MILLISECONDS.sleep(50);
                    serverResponse += "50 ms, ";
                    TimeUnit.MILLISECONDS.sleep(50);
                    serverResponse += "100 ms, ";
                    TimeUnit.MILLISECONDS.sleep(50);
                    serverResponse += "150 ms, ";
                    TimeUnit.MILLISECONDS.sleep(50);
                    serverResponse += "200 ms, ";
                    TimeUnit.MILLISECONDS.sleep(50);
                    serverResponse += "250 ms\n";
                    response = new DefaultFullHttpResponse(HTTP_1_1, OK, Unpooled.wrappedBuffer(serverResponse.getBytes(Charset.forName("UTF-8"))));
                    response.headers().set(CONTENT_LENGTH, response.content().readableBytes());
                    response.headers().set(CONNECTION, HttpHeaders.Values.KEEP_ALIVE);
                    ctx.write(response).addListener(ChannelFutureListener.CLOSE);
                } else {
                    FullHttpResponse response;
                    String serverResponse = String.format("NS Port: %d\n", port);
                    TimeUnit.MILLISECONDS.sleep(50);
                    serverResponse += "50 ms, ";
                    TimeUnit.MILLISECONDS.sleep(50);
                    serverResponse += "100 ms, ";
                    TimeUnit.MILLISECONDS.sleep(50);
                    serverResponse += "150 ms, ";
                    TimeUnit.MILLISECONDS.sleep(50);
                    serverResponse += "200 ms, ";
                    TimeUnit.MILLISECONDS.sleep(50);
                    serverResponse += "250 ms\n";
                    response = new DefaultFullHttpResponse(HTTP_1_1, OK, Unpooled.wrappedBuffer(serverResponse.getBytes(Charset.forName("UTF-8"))));
                    response.headers().set(CONTENT_LENGTH, response.content().readableBytes());
                    response.headers().set(CONNECTION, HttpHeaders.Values.KEEP_ALIVE);
                    ctx.write(response);
                }
            }
        }

        @Override
        public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) throws Exception {
            cause.printStackTrace();
            ctx.close();
        }
    }
}
