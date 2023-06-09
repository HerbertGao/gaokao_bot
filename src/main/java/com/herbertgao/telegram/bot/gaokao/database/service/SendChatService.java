package com.herbertgao.telegram.bot.gaokao.database.service;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.herbertgao.telegram.bot.gaokao.database.domain.SendChat;
import com.herbertgao.telegram.bot.gaokao.database.mapper.SendChatMapper;
import org.springframework.stereotype.Service;

@Service
public class SendChatService extends ServiceImpl<SendChatMapper, SendChat> {
}
