<?php
// GENERATED CODE -- DO NOT EDIT!

namespace Gufo;

/**
 * ============================================================
 * Gufo Core Transport (Dynamic Execution Bus)
 * ============================================================
 * - Keeps URL structure compatibility:
 *   /api/v1/{module}/{param}/{param_id}/{param_idd}
 * - Supports HTTP + gRPC
 * - Separates routing, query, and body
 * ============================================================
 *
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
     * @param \Gufo\Request $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function Do(\Gufo\Request $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/gufo.Reverse/Do',
        $argument,
        ['\Gufo\Response', 'decode'],
        $metadata, $options);
    }

    /**
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\BidiStreamingCall
     */
    public function Stream($metadata = [], $options = []) {
        return $this->_bidiRequest('/gufo.Reverse/Stream',
        ['\Gufo\Response','decode'],
        $metadata, $options);
    }

}
