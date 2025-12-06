<?php
// GENERATED CODE -- DO NOT EDIT!

namespace ;

/**
 * =============================
 * Gufo gRPC Service Definition
 * =============================
 */
class ReverseClient extends \Grpc\BaseStub {

    /**
     * @param string $hostname hostname
     * @param array $opts channel options
     * @param \Grpc\Channel $channel (optional) re-use channel object
     */
    public function __construct($hostname, $opts, $channel = null) {
        parent::__construct($hostname, $opts, $channel);
    }

    /**
     * @param \Request $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function Do(\Request $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/Reverse/Do',
        $argument,
        ['\Response', 'decode'],
        $metadata, $options);
    }

    /**
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\BidiStreamingCall
     */
    public function Stream($metadata = [], $options = []) {
        return $this->_bidiRequest('/Reverse/Stream',
        ['\Response','decode'],
        $metadata, $options);
    }

}
