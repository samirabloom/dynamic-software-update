import io.netty.bootstrap.ServerBootstrap;
import io.netty.buffer.Unpooled;
import io.netty.channel.*;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.SocketChannel;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import io.netty.handler.codec.http.DefaultFullHttpResponse;
import io.netty.handler.codec.http.FullHttpResponse;
import io.netty.handler.codec.http.HttpRequest;
import io.netty.handler.codec.http.HttpServerCodec;
import io.netty.handler.logging.LoggingHandler;

import java.nio.charset.Charset;

import static io.netty.handler.codec.http.HttpHeaders.Names.CONTENT_LENGTH;
import static io.netty.handler.codec.http.HttpHeaders.Names.CONTENT_TYPE;
import static io.netty.handler.codec.http.HttpResponseStatus.NOT_FOUND;
import static io.netty.handler.codec.http.HttpResponseStatus.OK;
import static io.netty.handler.codec.http.HttpVersion.HTTP_1_1;

public class TestServer {

    public static void main(String[] args) throws InterruptedException {
        System.out.println("STARTING SERVER FOR HTTP ON PORT: " + 8080);
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
                            pipeline.addLast("simple_test_handler", new TestServerHandler());
                        }
                    })
                    .bind(8080).channel().closeFuture().sync();
        } catch (Exception e) {
            throw new RuntimeException("Exception running test server", e);
        } finally {
            bossGroup.shutdownGracefully();
            workerGroup.shutdownGracefully();
        }
    }

    private static class TestServerHandler extends SimpleChannelInboundHandler<Object> {

        @Override
        public void channelReadComplete(ChannelHandlerContext ctx) {
            ctx.flush();
        }

        @Override
        public void channelRead0(ChannelHandlerContext ctx, Object msg) throws Exception {
            if (msg instanceof HttpRequest) {
                HttpRequest req = (HttpRequest) msg;

                FullHttpResponse response;
                if (req.getUri().equals("/unknown")) {
                    response = new DefaultFullHttpResponse(HTTP_1_1, NOT_FOUND);
                } else if (req.getUri().equals("/zeromq")) {
                    response = new DefaultFullHttpResponse(HTTP_1_1, OK, Unpooled.wrappedBuffer("it is fun to play with ZeroMQ".getBytes(Charset.forName("UTF-8"))));
                    response.headers().set("X-ZeroMQ", "a_header_value");
                    response.headers().set(CONTENT_TYPE, "text/plain");
                    response.headers().set(CONTENT_LENGTH, response.content().readableBytes());
                } else if (req.getUri().equals("/netty")) {
                    response = new DefaultFullHttpResponse(HTTP_1_1, OK, Unpooled.wrappedBuffer("it is also fun to play with Netty".getBytes(Charset.forName("UTF-8"))));
                    response.headers().set("X-Netty", "a_header_value");
                    response.headers().set(CONTENT_TYPE, "text/plain");
                    response.headers().set(CONTENT_LENGTH, response.content().readableBytes());
                } else {
                    response = new DefaultFullHttpResponse(HTTP_1_1, OK, Unpooled.wrappedBuffer("Hello World".getBytes(Charset.forName("UTF-8"))));
                    response.headers().set(CONTENT_TYPE, "text/plain");
                    response.headers().set(CONTENT_LENGTH, response.content().readableBytes());
                }
                ctx.write(response).addListener(ChannelFutureListener.CLOSE);
            }
        }

        @Override
        public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) throws Exception {
            cause.printStackTrace();
            ctx.close();
        }
    }
}